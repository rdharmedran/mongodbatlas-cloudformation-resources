// Code generated by 'cfn generate', changes will be undone by the next invocation. DO NOT EDIT.
package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/cloud-provider-snapshot-restore-jobs/cmd/resource"
	"go.mongodb.org/atlas"
)

// Handler is a container for the CRUDL actions exported by resources
type Handler struct{}

// Create wraps the related Create function exposed by the resource code
func (r *Handler) Create(req handler.Request) handler.ProgressEvent {
	return wrap(req, resource.Create)
}

// Read wraps the related Read function exposed by the resource code
func (r *Handler) Read(req handler.Request) handler.ProgressEvent {
	return wrap(req, resource.Read)
}

// Update wraps the related Update function exposed by the resource code
func (r *Handler) Update(req handler.Request) handler.ProgressEvent {
	return wrap(req, resource.Update)
}

// Delete wraps the related Delete function exposed by the resource code
func (r *Handler) Delete(req handler.Request) handler.ProgressEvent {
	return wrap(req, resource.Delete)
}

// List wraps the related List function exposed by the resource code
func (r *Handler) List(req handler.Request) handler.ProgressEvent {
	return wrap(req, resource.List)
}

// main is the entry point of the application.
func main() {
	cfn.Start(&Handler{})
}

type handlerFunc func(handler.Request, *resource.Model, *resource.Model) (handler.ProgressEvent, error)

func wrap(req handler.Request, f handlerFunc) (response handler.ProgressEvent) {
	defer func() {
		// Catch any panics and return a failed ProgressEvent
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = errors.New(fmt.Sprint(r))
			}

			log.Printf("Trapped error in handler: %v", err)

			response = handler.NewFailedEvent(err)
		}
	}()

	// Populate the previous model
	//prevModel := &resource.Model{}
	currentModel := &atlas.CloudProviderSnapshotRestoreJobs
	// Here we will tell the tool WHAT type we want, we already have our

	// ************
	// {'models': {'ResourceModel': {'ProjectId': <ResolvedType(ContainerType.PRIMITIVE, string)>, 'ClusterName': <ResolvedType(ContainerType.PRIMITIVE, string)>, 'Id': <ResolvedType(ContainerType.PRIMITIVE, string)>, 'DeliveryType': <ResolvedType(ContainerType.PRIMITIVE, string)>, 'DeliveryUrl': <ResolvedType(ContainerType.LIST, <ResolvedType(ContainerType.PRIMITIVE, string)>)>, 'Cancelled': <ResolvedType(ContainerType.PRIMITIVE, boolean)>, 'CreatedAt': <ResolvedType(ContainerType.PRIMITIVE, string)>, 'Expired': <ResolvedType(ContainerType.PRIMITIVE, boolean)>, 'ExpiresAt': <ResolvedType(ContainerType.PRIMITIVE, string)>, 'FinishedAt': <ResolvedType(ContainerType.PRIMITIVE, string)>, 'Timestamp': <ResolvedType(ContainerType.PRIMITIVE, string)>, 'SnapshotId': <ResolvedType(ContainerType.PRIMITIVE, string)>, 'Links': <ResolvedType(ContainerType.LIST, <ResolvedType(ContainerType.MODEL, Links)>)>, 'OpLogTs': <ResolvedType(ContainerType.PRIMITIVE, string)>, 'PointInTimeUtcSeconds': <ResolvedType(ContainerType.PRIMITIVE, integer)>, 'TargetProjectId': <ResolvedType(ContainerType.PRIMITIVE, string)>, 'TargetClusterName': <ResolvedType(ContainerType.PRIMITIVE, string)>, 'ApiKeys': <ResolvedType(ContainerType.MODEL, ApiKeyDefinition)>}, 'Links': {'Rel': <ResolvedType(ContainerType.PRIMITIVE, string)>, 'Href': <ResolvedType(ContainerType.PRIMITIVE, string)>}, 'ApiKeyDefinition': {'PublicKey': <ResolvedType(ContainerType.PRIMITIVE, string)>, 'PrivateKey': <ResolvedType(ContainerType.PRIMITIVE, string)>}}, 'path': PosixPath('/home/jason/work/mongodbatlas-cloudformation-resources/cfn-resources/cloud-provider-snapshot-restore-jobs/cmd/main.go'), 'type_name': 'MongoDB::Atlas::CloudProviderSnapshotRestoreJobs', 'client_type': 'CloudProviderSnapshotRestoreJobs', 'imports': 'go.mongodb.org/atlas'}
	// ************

	// types in the mongodbatlas go client!
	if err := req.UnmarshalPrevious(prevModel); err != nil {
		log.Printf("Error unmarshaling prev model: %v", err)
		return handler.NewFailedEvent(err)
	}

	// Populate the current model
	//currentModel := &resource.Model{}
	currentModel := &atlas.CloudProviderSnapshotRestoreJobs

	if err := req.Unmarshal(currentModel); err != nil {
		log.Printf("Error unmarshaling model: %v", err)
		return handler.NewFailedEvent(err)
	}

	response, err := f(req, prevModel, currentModel)
	if err != nil {
		log.Printf("Error returned from handler function: %v", err)
		return handler.NewFailedEvent(err)
	}

	return response
}
