# flowdock-notification-resource
Resource to send notifications from [Concourse](https://concourse-ci.org/) to a flow on [Flowdock](https://flowdock.com)

NOTE: This resource uses the [Messages API](https://www.flowdock.com/api/messages) from Flowdock and relies on `flow_token` for [authentication](https://www.flowdock.com/api/authentication#source-token)

As a pre-requisite, you should configure a [source](https://www.flowdock.com/api/sources) for a flow. The source will then have a `flow_token` which will allow concourse to post messages to that flow.

## Usage
Include the following in your Pipeline YAML file, replacing the values in the parentheses `(())`
```
resource_types:
- name: flowdock-notification
  type: docker-image
  source:
    repository: rafikn/flowdock-notification-resource

resources:
- name: flowdock-notifier
  type: flowdock-notification
  source:
    author: concourse
    avatar: http://cl.ly/image/3e1h0H3H2s0P/concourse-logo.png
    flow_token: ((flow_token))
    flow_api: ((flow_api))
```
Then configure your job to use this resource
```
- name: job
  plan:
  - task: sample_task
    ...
    on_success:
      put: flowdock-notifier
      params:
        event: activity
        message_title: sample_task job success
        message_body: sample_task job has succeeded
        status_colour: lime
        status_value: SUCCESS
    on_failure:
      put: flowdock-notifier
      params:
        event: activity
        message_title: sample_task job failure
        message_body: sample_task job has failed
        status_colour: red
        status_value: FAILURE
```
## Source Configuration
* `flow_token` Required. See [this](https://www.flowdock.com/api/authentication#source-token) on how to get a `flow_token`. A `flow_token` is linked to a [Source](https://www.flowdock.com/api/sources). This 'intergration' allows concourse to post messages to that flow.
* `flow_api` Required. Can be `https://api.flowdock.com/messages`, `https://api.flowdock.com/flows/:organization/:flow/messages` or `https://api.flowdock.com/flows/:organization/:flow/messages/threads/:id/messages`
* `author` Username associated with the notification
* `avatar` Avatar associated with the notification
* `event` Can be `activity` or `message`. `activity` notifications are shown in the Inbox of the flow. `message` notification appear in the main thread

## `out` Configuration
* `message_body` Required. Text of the message.
* `message_title` Optional. (`activity` notification only) Title for the message. If provided, the build number will be appended to it, if none provided, the default pattern is `pipeline | job | build | [status]`
* `status_value` Required. (`activity` notification only) Status label for the thread
* `status_color` Required. (`activity` notification only) Background colour for the status label
* `thread_id` Optional. Notifications are threaded by job's name by default or if set to `job_name`. Thread by exact build number if set to `build_number`. Otherwise, thread ID is explicitly set to the configured value of this param.
* `version_file` Optional. Logical path to a version file as the concourse semver resource would expect. `message_title` and `message_body` can then include the text `%version%`, which will be substituted by the version number at runtime. Useful for application release messages.

All configuration of the resource's `source` can be overridden in the notification task's `params`

Typically, it's best to configure common values, e.g. `flow_token`, `flow_api` etc... at the `source` level and configure the message properties, e.g. `event`, `message_body` etc... at the tasks's `params` level

## Contributing
Please make all pull requests to the master branch and ensure the `out` script tests pass locally.

## License
MIT
