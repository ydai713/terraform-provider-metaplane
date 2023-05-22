package provider

import (
  "context"
  "os"
  
  "github.com/klaviyo/terraform-provider-metaplane/internal/api"
  
  "github.com/hashicorp/terraform-plugin-framework/path"
  "github.com/hashicorp/terraform-plugin-framework/datasource"
  "github.com/hashicorp/terraform-plugin-framework/provider"
  "github.com/hashicorp/terraform-plugin-framework/provider/schema"
  "github.com/hashicorp/terraform-plugin-framework/resource"
  "github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure metaplaneProvider satisfies various provider interfaces.
var _ provider.Provider = &metaplaneProvider{}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &metaplaneProvider{
			version: version,
		}
	}
}

// metaplaneProvider defines the provider implementation.
type metaplaneProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// metaplaneProviderModel describes the provider data model.
type metaplaneProviderModel struct {
	ApiKey types.String `tfsdk:"api_key"`
}

func (p *metaplaneProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "metaplane"
	resp.Version = p.version
}

func (p *metaplaneProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Metaplane API Key",
				Optional:            true,
			},
		},
	}
}

func (p *metaplaneProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
  // Retrieve provider data from configuration
    var config metaplaneProviderModel
    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // If practitioner provided a configuration value for any of the
    // attributes, it must be a known value.

    if config.ApiKey.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("api_key"),
            "Unknown Metaplane API Key",
            "The provider cannot create the Metaplane API client as there is an unknown configuration value for the Metaplane API Key. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the METAPLANE_API_KEY environment variable.",
        )
    }

    if resp.Diagnostics.HasError() {
        return
    }

    // Default values to environment variables, but override
    // with Terraform configuration value if set.
    api_key := os.Getenv("METAPLANE_API_KEY")

    if !config.ApiKey.IsNull() {
        api_key = config.ApiKey.ValueString()
    }

    // If any of the expected configurations are missing, return
    // errors with provider-specific guidance.
    if api_key == "" {
        resp.Diagnostics.AddAttributeError(
            path.Root("api_key"),
            "Missing Metaplane API Host",
            "The provider cannot create the Metaplane API client as there is a missing or empty value for the Metaplane API kye. "+
                "Set the api_key value in the configuration or use the METAPLANE_API_KEY environment variable. "+
                "If either is already set, ensure the value is not empty.",
        )
    }


    if resp.Diagnostics.HasError() {
        return
    }

    // Create a new metaplane client using the configuration values
    client := api.NewClient(&api_key)

    // Make the metaplane client available during DataSource and Resource
    // type Configure methods.
    resp.DataSourceData = client
    resp.ResourceData = client
}

func (p *metaplaneProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
    NewMonitorResource,
	}
}

func (p *metaplaneProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewMonitorDataSource,
		NewConnectionDataSource,
	}
}

