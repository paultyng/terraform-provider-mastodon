package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mattn/go-mastodon"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &FollowResource{}
var _ resource.ResourceWithImportState = &FollowResource{}

func NewFollowResource() resource.Resource {
	return &FollowResource{}
}

// FollowResource defines the resource implementation.
type FollowResource struct {
	client *mastodon.Client
}

// FollowResourceModel describes the resource data model.
type FollowResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Account types.String `tfsdk:"account"`
	// Reblogs types.Bool   `tfsdk:"reblogs"`
	// Notify  types.Bool   `tfsdk:"notify"`
}

func (r *FollowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_follow"
}

func (r *FollowResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "mastodon_follow manages the following relationship between the user and an account, either local or remote.",

		Attributes: map[string]tfsdk.Attribute{
			"account": {
				Required:            true,
				MarkdownDescription: "The name of the account to follow. If remote, it must include the domain.",
				Type:                types.StringType,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},

			// "reblogs": {
			// 	// Optional:            true,
			// 	Computed:            true,
			// 	MarkdownDescription: "",
			// 	Type:                types.BoolType,
			// },
			// "notify": {
			// 	// Optional:            true,
			// 	Computed:            true,
			// 	MarkdownDescription: "",
			// 	Type:                types.BoolType,
			// },

			"id": {
				Computed:            true,
				MarkdownDescription: "The ID of the account on the local server.",
				Type:                types.StringType,
			},
		},
	}, nil
}

func (r *FollowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*mastodon.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *FollowResource) getAccountID(ctx context.Context, account string) (mastodon.ID, error) {
	accounts, err := r.client.AccountsSearch(ctx, account, 1)
	if err != nil {
		return "", fmt.Errorf("unable to search for account: %w", err)
	}
	if len(accounts) != 1 {
		return "", fmt.Errorf("unable to find exact match, found %d", len(accounts))
	}
	return accounts[0].ID, nil
}

func (r *FollowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *FollowResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	accountID, err := r.getAccountID(ctx, data.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to find account, got error: %s", err))
		return
	}

	_, err = r.client.AccountFollow(ctx, accountID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to follow account, got error: %s", err))
		return
	}

	data.ID = types.StringValue(string(accountID))
	// data.Reblogs = types.BoolValue(relationship.ShowingReblogs)
	// data.Notify = types.BoolValue(relationship.Notifying)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FollowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *FollowResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	relationships, err := r.client.GetAccountRelationships(ctx, []string{data.ID.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get account relationship, got error: %s", err))
		return
	}
	if len(relationships) > 1 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to find relationship, found %d", len(relationships)))
		return
	}
	if len(relationships) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FollowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *FollowResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	relationships, err := r.client.GetAccountRelationships(ctx, []string{data.ID.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get account relationship, got error: %s", err))
		return
	}
	if len(relationships) != 1 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to find exact relationship, found %d", len(relationships)))
		return
	}

	// nothing to do here yet

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FollowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *FollowResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	relationships, err := r.client.GetAccountRelationships(ctx, []string{data.ID.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get account relationship, got error: %s", err))
		return
	}
	if len(relationships) > 1 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to find relationship, found %d", len(relationships)))
		return
	}
	if len(relationships) == 0 {
		// nothing to do, already gone
		return
	}

	_, err = r.client.AccountUnfollow(ctx, mastodon.ID(data.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to unfollow account, got error: %s", err))
		return
	}
}

func (r *FollowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	accountID, err := r.getAccountID(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to find account, got error: %s", err))
		return
	}

	data := &FollowResourceModel{
		ID:      types.StringValue(string(accountID)),
		Account: types.StringValue(req.ID),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
