package datasets

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"govchain/x/datasets/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListEntry",
					Use:       "list-entry",
					Short:     "List all entry",
				},
				{
					RpcMethod:      "GetEntry",
					Use:            "get-entry [id]",
					Short:          "Gets a entry by id",
					Alias:          []string{"show-entry"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod:      "EntriesByAgency",
					Use:            "entries-by-agency [agency]",
					Short:          "Query entries-by-agency",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "agency"}},
				},

				{
					RpcMethod:      "EntriesByCategory",
					Use:            "entries-by-category [category]",
					Short:          "Query entries-by-category",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "category"}},
				},

				{
					RpcMethod:      "EntriesByMimetype",
					Use:            "entries-by-mimetype [mime-type]",
					Short:          "Query entries-by-mimetype",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "mime_type"}},
				},

				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateEntry",
					Use:            "create-entry [title] [description] [ipfs-cid] [mime-type] [file-name] [file-url] [fallback-url] [file-size] [checksum-sha-256] [agency] [category] [submitter] [timestamp] [pin-count]",
					Short:          "Create entry",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "title"}, {ProtoField: "description"}, {ProtoField: "ipfs_cid"}, {ProtoField: "mime_type"}, {ProtoField: "file_name"}, {ProtoField: "file_url"}, {ProtoField: "fallback_url"}, {ProtoField: "file_size"}, {ProtoField: "checksum_sha_256"}, {ProtoField: "agency"}, {ProtoField: "category"}, {ProtoField: "submitter"}, {ProtoField: "timestamp"}, {ProtoField: "pin_count"}},
				},
				{
					RpcMethod:      "UpdateEntry",
					Use:            "update-entry [id] [title] [description] [ipfs-cid] [mime-type] [file-name] [file-url] [fallback-url] [file-size] [checksum-sha-256] [agency] [category] [submitter] [timestamp] [pin-count]",
					Short:          "Update entry",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "title"}, {ProtoField: "description"}, {ProtoField: "ipfs_cid"}, {ProtoField: "mime_type"}, {ProtoField: "file_name"}, {ProtoField: "file_url"}, {ProtoField: "fallback_url"}, {ProtoField: "file_size"}, {ProtoField: "checksum_sha_256"}, {ProtoField: "agency"}, {ProtoField: "category"}, {ProtoField: "submitter"}, {ProtoField: "timestamp"}, {ProtoField: "pin_count"}},
				},
				{
					RpcMethod:      "DeleteEntry",
					Use:            "delete-entry [id]",
					Short:          "Delete entry",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
