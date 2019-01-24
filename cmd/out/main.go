package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/blang/semver"
)

func metadataMap() []map[string]string {
	metadataMap := make([]map[string]string, 6)

	metadataMap[0] = buildField("Team", os.Getenv(BuildTeamName))
	metadataMap[1] = buildField("Pipeline", os.Getenv(BuildPipelineName))
	metadataMap[2] = buildField("Job", os.Getenv(BuildJobName))
	metadataMap[3] = buildField("Build Number", os.Getenv(BuildName))
	metadataMap[4] = buildField("Build ID", os.Getenv(BuildId))
	metadataMap[5] = buildField("Concourse URL", os.Getenv(AtcExternalUrl))

	return metadataMap
}

func externalUrl() string {
	return fmt.Sprintf("%s/teams/%s/pipelines/%s/jobs/%s/builds/%s",
		os.Getenv(AtcExternalUrl),
		os.Getenv(BuildTeamName),
		os.Getenv(BuildPipelineName),
		os.Getenv(BuildJobName),
		os.Getenv(BuildName),
	)
}

func buildField(label string, value string) map[string]string {
	return map[string]string{
		"label": label,
		"value": value,
	}
}

func readInput(input []byte, source *Input) {
	err := json.Unmarshal(input, &source)
	if err != nil {
		os.Stderr.Write([]byte(fmt.Sprintf("error parsing input string to json %s\n", err)))
	}
}

func buildRequestData(config *Input) map[string]interface{} {

	source := config.Source
	params := config.Params
	metadata := metadataMap()

	flowToken := source.FlowToken
	if params.FlowToken != "" {
		flowToken = params.FlowToken
	}

	event := source.Event
	if params.Event != "" {
		event = params.Event
	}

	authorName := source.Author
	if params.Author != "" {
		authorName = params.Author
	}

	authorAvatar := source.Avatar
	if params.Author != "" {
		authorAvatar = params.Avatar
	}

	team := os.Getenv(BuildTeamName)
	pipeline := os.Getenv(BuildPipelineName)
	job := os.Getenv(BuildJobName)
	build := os.Getenv(BuildName)

	threadId := source.ThreadId
	if params.ThreadId != "" {
		threadId = params.ThreadId
	}
	externalThreadId := threadId
	switch threadId {
	case "", "job_name":
		externalThreadId = fmt.Sprintf("%s_%s_%s", team, pipeline, job)
	case "build_number":
		externalThreadId = fmt.Sprintf("%s_%s_%s_%s", team, pipeline, job, build)
	}

	threadTitle := fmt.Sprintf("%s | %s | %s", pipeline, job, build)

	eventTitle := params.MessageTitle
	if eventTitle == "" {
		eventTitle = fmt.Sprintf("%s | %s | %s [%s]", pipeline, job, build, params.StatusValue)
	} else {
		eventTitle = fmt.Sprintf("%s #%s", eventTitle, build)
	}

	messageBody := params.MessageBody // only used when event == message

	if params.VersionFile != "" {
		workdir := os.Args[1] // per concourse spec first arg is the target dir
		versionFilePath := fmt.Sprintf("%s/%s", workdir, params.VersionFile)

		version, err := ioutil.ReadFile(versionFilePath)
		if err != nil {
			panic(err)
		}

		_, err = semver.Parse(fmt.Sprintf("%s", version)) // parse to validate the semver
		if err != nil {
			panic(err)
		}

		// TODO: If %version% present in either of MessageBody or MessageTitle replace
		if strings.Contains(eventTitle, `%version%`) {
			eventTitle = strings.Replace(eventTitle, `%version%`, fmt.Sprintf("%s", version), -1)
		}

		if strings.Contains(messageBody, `%version%`) {
			messageBody = strings.Replace(messageBody, `%version%`, fmt.Sprintf("%s", version), -1)
		}
	}

	jsonData := map[string]interface{}{
		"flow_token": flowToken,
		"event":      event,
		"content":    messageBody, // only used when event == message
		"author": map[string]string{
			"name":   authorName,
			"avatar": authorAvatar,
		},
		"title":              eventTitle,
		"external_thread_id": externalThreadId,
		"thread": map[string]interface{}{
			"title": threadTitle,
			"body":  messageBody,
			"status": map[string]string{
				"color": params.StatusColour,
				"value": params.StatusValue,
			},
			"fields":       metadata,
			"external_url": externalUrl(),
		},
	}

	return jsonData
}

func sendRequest(requestUrl string, requestData map[string]interface{}) {
	jsonString, err := json.Marshal(requestData)
	if err != nil {
		os.Stderr.Write([]byte(fmt.Sprintf("error parsing request json body to string %s\n", err)))
	}

	result, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(jsonString))
	if err != nil {
		os.Stderr.Write([]byte(fmt.Sprintf("The HTTP request failed with error %s\n", err)))
	}

	if int(result.StatusCode) >= 300 {
		os.Stderr.Write([]byte("Result was not ok"))
		os.Stderr.Write([]byte(fmt.Sprintln(result)))
	} else {
		os.Stderr.Write([]byte("Successfully notified flow"))
	}
}

func main() {
	scanner := bufio.NewReader(os.Stdin)
	line, _, err := scanner.ReadLine()
	if err != nil {
		os.Stderr.Write([]byte("error reading input\n"))
		scanner = bufio.NewReader(os.Stdin)
		scanner.WriteTo(os.Stderr)
	}

	var input Input
	readInput(line, &input)

	requestData := buildRequestData(&input)
	requestUrl := input.Source.FlowApi
	if input.Params.FlowApi != "" {
		requestUrl = input.Params.FlowApi
	}

	sendRequest(requestUrl, requestData)
	os.Stdout.Write([]byte(fmt.Sprintf("{ \"version\" :{ \"ref\" :\"%s\"}}", input.Params.StatusValue)))
}

type Resource struct {
	FlowToken    string `json:"flow_token"`
	FlowApi      string `json:"flow_api"`
	Event        string `json:"event"`
	Author       string `json:"author"`
	Avatar       string `json:"avatar"`
	MessageTitle string `json:"message_title"`
	MessageBody  string `json:"message_body"`
	StatusColour string `json:"status_colour"`
	StatusValue  string `json:"status_value"`
	ThreadId     string `json:"thread_id"`
	VersionFile  string `json:"version_file"` //e.g. a version file as required by the semver resource.
}

type Input struct {
	Source *Resource `json:"source"`
	Params *Resource `json:"params"`
}

var BuildId = "BUILD_ID"
var BuildName = "BUILD_NAME"
var BuildJobName = "BUILD_JOB_NAME"
var BuildPipelineName = "BUILD_PIPELINE_NAME"
var BuildTeamName = "BUILD_TEAM_NAME"
var AtcExternalUrl = "ATC_EXTERNAL_URL"
