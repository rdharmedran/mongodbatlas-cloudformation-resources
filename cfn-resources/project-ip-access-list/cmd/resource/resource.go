package resource

import (
	"context"
    "strings"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util"
	"github.com/rs/xid"
	"go.mongodb.org/atlas/mongodbatlas"
    "github.com/aws/aws-sdk-go/service/cloudformation"
    log "github.com/sirupsen/logrus"
)

func init() {
    util.InitLogger()
}

// Create handles the Create event from the Cloudformation service.
func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
    log.Info("Create handler called")
	log.Debugf("currentModel: %+v, prevModel: %+v", currentModel, prevModel)
    
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		return handler.ProgressEvent{
            OperationStatus: handler.Failed,
            Message: err.Error(),
            HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil
	}

    event, err := createEntries(currentModel, client)
	if err != nil {
	    log.Errorf("Create err:%v", err)
		return event, nil
	}

	guid := xid.New()

	x := guid.String()
	currentModel.Id = &x
    log.WithFields(log.Fields{"currentModel":currentModel,}).Debug("+++ Create() +++++++++++++")
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Create Complete",
		ResourceModel:   currentModel,
	}, nil
}

// Read handles the Read event from the Cloudformation service.
func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		return handler.ProgressEvent{
            OperationStatus: handler.Failed,
            Message: err.Error(),
            HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil
	}

	projectID := *currentModel.ProjectId

    log.WithFields(log.Fields{"currentModel":currentModel,}).Debug("read read read         Read()     read read read")
    
	entries := []string{}
	for i, _ := range currentModel.AccessList {
        wl := currentModel.AccessList[i]
		entry := getEntry(wl)
		entries = append(entries, entry)
	}

    log.WithFields(log.Fields{"entries":entries,}).Debug("read read read         Read()     read read read")
	accesslist, progressEvent, err := getProjectIPAccessList(projectID, entries, client)
    log.WithFields(log.Fields{"accesslist":accesslist,"err":err,"progressEvent":progressEvent,}).Debug("read read read         Read()     read read read")
	if err != nil {
        log.Errorf("error READ access list projectID:%s, error: %s, progressEvent: %+v", projectID,err, progressEvent)
        return progressEvent, nil
	}

    
	currentModel.AccessList = flattenAccessList(currentModel.AccessList,accesslist)
    log.WithFields(log.Fields{">>>>>>> currentModel.Accesslist":currentModel.AccessList,}).Debug("Read()")

	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Read Complete",
		ResourceModel:   currentModel,
	}, nil

}

// Update handles the Update event from the Cloudformation service.
func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		return handler.ProgressEvent{
            OperationStatus: handler.Failed,
            Message: err.Error(),
            HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil
	}

    progressEvent, err := deleteEntries(currentModel, client)
	if err != nil {
        log.WithFields(log.Fields{"err":err,}).Error("deleteEntries error")
        return progressEvent, nil
	}

	progressEvent, err = createEntries(currentModel, client)
	if err != nil {
        log.WithFields(log.Fields{"err":err,}).Error("createEntries error")
		return progressEvent, nil
	}

	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Update Complete",
		ResourceModel:   currentModel,
	}, nil

}

// Delete handles the Delete event from the Cloudformation service.
func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		return handler.ProgressEvent{
            OperationStatus: handler.Failed,
            Message: err.Error(),
            HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil
	}

    event, err := deleteEntries(currentModel, client)
	if err != nil {
        log.WithFields(log.Fields{"err":err,"event":event,}).Error("Delete() deleteEntries error")
		return event, nil
	}

	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Delete Complete",
	}, nil

}

// List handles the List event from the Cloudformation service.
// NO-OP
func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	log.Printf("Got list request - returning read - %v", currentModel)
    readEvent, err := Read(req, prevModel, currentModel)
    log.Debugf("List readEvent:+%v   --------------------------- error:%+v",readEvent,err)
    if readEvent.OperationStatus==handler.Failed {
        return readEvent, nil
    }
    // create list with 1
    models := []interface{} {}
	models = append(models, readEvent.ResourceModel)

	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "List Complete",
		ResourceModels:   models,
	}, nil

}

func getProjectIPAccessList(projectID string, entries []string, conn *mongodbatlas.Client) ([]*mongodbatlas.ProjectIPAccessList, handler.ProgressEvent, error) {

	var accesslist []*mongodbatlas.ProjectIPAccessList
	for i, _ := range entries {
        entry := entries[i]
		result, resp, err := conn.ProjectIPAccessList.Get(context.Background(), projectID, entry)
		if err != nil {
            if resp != nil && resp.StatusCode == 404 {
                log.Warnf("Resource Not Found 404 for READ projectId:%s, entry:%+v, err:%+v", projectID, entry, err)
                return nil, handler.ProgressEvent{
                    Message: err.Error(),
                    OperationStatus: handler.Failed,
                    HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, err
            } else {
                log.Errorf("Error READ projectId:%s, err:%+v", projectID, err)
                return nil, handler.ProgressEvent{
                    Message: err.Error(),
                    OperationStatus: handler.Failed,
                    HandlerErrorCode: cloudformation.HandlerErrorCodeServiceInternalError}, err
            }
		}
        log.Debugf("%+v",strings.Split(result.CIDRBlock,"/"))
        // if "/32" cidr "it's just you" then convert to the ip addres field
        //if len(strings.Split(entry,"/")==0 {    //  an IPAddress, get's returned 
        /*
        ss := strings.Split(result.CIDRBlock,"/")
        if ss[1] == "32" {
            result.IPAddress = ss[0]
            result.CIDRBlock = ""
        }
        */
        log.WithFields(log.Fields{"result":result,}).Debug("getProjectIPAccessList")
		accesslist = append(accesslist, result)
	}
	return accesslist, handler.ProgressEvent{}, nil
}

func getProjectIPAccessListRequest(model *Model) []*mongodbatlas.ProjectIPAccessList {
	var accesslist []*mongodbatlas.ProjectIPAccessList
	for i, _ := range model.AccessList {
        w := model.AccessList[i]
		wl := &mongodbatlas.ProjectIPAccessList{}
		if w.Comment != nil {
			wl.Comment = *w.Comment
		}
		if w.CIDRBlock != nil {
			wl.CIDRBlock = *w.CIDRBlock
		}
		if w.IPAddress != nil {
			wl.IPAddress = *w.IPAddress
		}
		if w.AwsSecurityGroup != nil {
			wl.AwsSecurityGroup = *w.AwsSecurityGroup
		}

		log.Debugf("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ getProjectIPAccessListRequest: %+v\n", wl)

		accesslist = append(accesslist, wl)
	}
	log.Printf("getProjectIPAccessListRequest accesslist:%v", accesslist)
	return accesslist
}

func getEntry(wl AccessListDefinition) string {

	if wl.CIDRBlock != nil {
		return *wl.CIDRBlock
	}
	if wl.AwsSecurityGroup != nil {
		return *wl.AwsSecurityGroup
	}
	if wl.IPAddress != nil {
		return *wl.IPAddress
	}
	return ""
}

func flattenAccessList(original []AccessListDefinition, accesslist []*mongodbatlas.ProjectIPAccessList) []AccessListDefinition {
	var results []AccessListDefinition
	for i, _ := range accesslist {
        wl := accesslist[i]
        // only add properties which were in model to begin with
		r := AccessListDefinition{
			IPAddress:        &wl.IPAddress,
			Comment:          &wl.Comment,
		}
        if original[i].CIDRBlock!=nil {
            r.CIDRBlock = &wl.CIDRBlock
        }
        if original[i].ProjectId!=nil {
            r.ProjectId =  &wl.GroupID
        }
        if original[i].AwsSecurityGroup!=nil {
            r.AwsSecurityGroup =  &wl.AwsSecurityGroup
        }
		results = append(results, r)
	}
	return results
}

func createEntries(model *Model, client *mongodbatlas.Client) (handler.ProgressEvent, error) {
	request := getProjectIPAccessListRequest(model)
	projectID := *model.ProjectId
    log.WithFields(log.Fields{"request":request,}).Warningf("createEntries : projectID:%s, model:%+v, request:%+v", projectID)
	result, _, err := client.ProjectIPAccessList.Create(context.Background(), projectID, request)
    if err != nil {
        log.WithFields(log.Fields{"request":request,}).Errorf("Error createEntries projectId:%s,err:%+v", projectID, err)
        return handler.ProgressEvent{
            Message: err.Error(),
            OperationStatus: handler.Failed,
            HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, err
    }
    log.WithFields(log.Fields{"result":result,}).Debug("createEntries")
	return handler.ProgressEvent{}, nil
}

func deleteEntries(model *Model, client *mongodbatlas.Client) (handler.ProgressEvent, error) {
	projectID := *model.ProjectId

	for i, _ := range model.AccessList {
        wl := model.AccessList[i]
		entry := getEntry(wl)
        log.WithFields(log.Fields{"entry":entry,}).Debug("deleteEntriesi 1")
		resp, errDelete := client.ProjectIPAccessList.Delete(context.Background(), projectID, entry)
        log.WithFields(log.Fields{"errDelete":errDelete,}).Debug("deleteEntries")
		if errDelete != nil {
            if resp != nil && resp.StatusCode == 404 {
                log.Warnf("Resource Not Found 404 deleteEntries projectId:%s, entry:%+v, err:%+v", projectID, entry, errDelete)
                return handler.ProgressEvent{
                    Message: errDelete.Error(),
                    OperationStatus: handler.Failed,
                    HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, errDelete
            } else {
                log.Errorf("Error READ projectId:%s, err:%+v", projectID, errDelete)
                return handler.ProgressEvent{
                    Message: errDelete.Error(),
                    OperationStatus: handler.Failed,
                    HandlerErrorCode: cloudformation.HandlerErrorCodeServiceInternalError}, errDelete
            }
		}
	}

	return handler.ProgressEvent{}, nil
}
