# Contributing

## Common Contributions

### Add a new resource

Example PR: [link](https://github.com/SumoLogic/terraform-provider-sumologic/pull/710)

* Code
    * `sumologic/sumologic_foo.go`
        * data type
        * functions to interact with [Sumo Logic API](https://help.sumologic.com/docs/api/)
    * `sumologic/resource_sumologic_foo.go`
        * schema
        * functions to interact with tf state
    * `sumologic/resource_sumologic_foo_test.go`
        * acceptance tests
    * `sumologic/provider.go`
        * add the new resource to the `ResourcesMap`
* Documentation ([preview tool](https://registry.terraform.io/tools/doc-preview))
    * `website/docs/r/foo.html.markdown`
        * purpose
        * example usage
        * detailed reference
    * `CHANGELOG.md`
        * leave the `## X.Y.Z (Unreleased)` header
            * a new version number will be selected only when a release is cut
        * add your new changes under `* Add new change notes here` (and do not delete this line)

### Add a new data source

Example PR: [link](https://github.com/SumoLogic/terraform-provider-sumologic/pull/762)

* Code
    * `sumologic/sumologic_foo.go`
        * data type
        * functions to interact with [Sumo Logic API](https://help.sumologic.com/docs/api/)
    * `sumologic/data_source_sumologic_foo.go` (singular, for reading one)
        * schema
        * functions to read into the tf state
    * `sumologic/data_source_sumologic_foos.go` (plural, for reading multiple)
        * schema
        * functions to read into the tf state
    * `sumologic/resource_sumologic_foo_test.go` (singular)
        * acceptance tests
    * `sumologic/resource_sumologic_foos_test.go` (plural)
        * acceptance tests
    * update `sumologic/provider.go`
        * add the new data sources to the `DataSourcesMap`
* Documentation ([preview tool](https://registry.terraform.io/tools/doc-preview))
    * `website/docs/d/foo.html.markdown` (singular)
        * purpose
        * example usage
        * detailed reference
    * `website/docs/d/foos.html.markdown` (plural)
        * purpose
        * example usage
        * detailed reference
    * `CHANGELOG.md`
        * leave the `## X.Y.Z (Unreleased)` header
            * a new version number will be selected only when a release is cut
        * add your new changes under `* Add new change notes here` (and do not delete this line)

If only the singular or plural versions of the data source are required, that's okay. It's not always necessary to have both.

### Cut a new release (Sumo Logic team members only)

Example PR: [link](https://github.com/SumoLogic/terraform-provider-sumologic/pull/759)

1. Code
    * `CHANGELOG.md`
        * Choose a new version number
        * Add a new line, under `* Add new change notes here`, with format `## MAJOR.MINOR.PATCH (Month Day, Year)`
        * Tidy up any formatting inconsistencies
2. Tag the commit
    ```bash
    git pull origin master
    git tag vX.XX.XX
    git push origin vX.XX.XX
    # where X.XX.XX is the new version of the provider
    ```
3. Wait for the goreleaser job
    * When you push the tag, it will trigger a GitHub Action job named `goreleaser`, which does the following automatically:
        * Create a draft release for the new tag
        * Build binaries
        * Attach those binaries to the draft release
    * You can find the job through the [Actions page](https://github.com/SumoLogic/terraform-provider-sumologic/actions)
        * [Example](https://github.com/SumoLogic/terraform-provider-sumologic/actions/runs/15145892343)
    * Wait for the job to finish (less than 10 minutes)
4. Publish the release
    * Open the [Releases page](https://github.com/SumoLogic/terraform-provider-sumologic/releases)
    * Click the `Edit` button for the new draft release (at the top right, looks like a pencil)
    * Copy and paste the release notes from `CHANGELOG.md` into the description
    * Scroll all the way down and ensure `Set as the latest release` is checked
    * Click `Publish release`
5. Verify that the release is published to the Terraform Registry
    * Open the [Terraform Registry](https://registry.terraform.io/providers/SumoLogic/sumologic/latest)
    * Check that the version number matches the new release
        * NOTE: This can sometimes take up to 30 minutes
    * If the registry still has not updated after a long time, you can redeliver the webhook messages:
        * Open the [Webhook settings](https://github.com/SumoLogic/terraform-provider-sumologic/settings/hooks)
        * For the Terraform Registry, click `Edit`
        * Select the `Recent Deliveries` tab
        * Select each of the recent messages and click `Redeliver` in order

## PRs from the Community

When a PR includes changes to the provider code, then a GitHub Action will attempt to run the acceptance tests.
The acceptance tests depend on several API keys which are stored in the GitHub Secrets. For safety and security,
these Secrets are not available in PRs from forked repositories. At the time of writing, our best workaround is
for a Sumo Logic member to duplicate the PR.

1. Review the PR changes
    * Correctness?
    * Completeness?
    * Security risks?
2. Duplicate the branch
    ```bash
    gh pr checkout PR_NUMBER
    git checkout -b BRANCH_NAME
    git push origin BRANCH_NAME
    ```
3. Create a new PR ([Example](https://github.com/SumoLogic/terraform-provider-sumologic/pull/797))
