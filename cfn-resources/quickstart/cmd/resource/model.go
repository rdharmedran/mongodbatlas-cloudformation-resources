// Code generated by 'cfn generate', changes will be undone by the next invocation. DO NOT EDIT.
// Updates to this type are made my editing the schema file and executing the 'generate' command.
package resource

// Model is autogenerated from the json schema
type Model struct {
	AtlasApiKeySecret *string `json:",omitempty"`
	ClusterName       *string `json:",omitempty"`
	Region            *string `json:",omitempty"`
	ClusterSize       *string `json:",omitempty"`
	AddSampleData     *bool   `json:",omitempty"`
	IAMRole           *string `json:",omitempty"`
	ConnectionString  *string `json:",omitempty"`
}