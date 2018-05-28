package main

import (
	"testing"
	"strings"
	"os"
	"fmt"
)

func TestReadInput(t *testing.T) {

	sample := `{"source":{"author":"concourse","avatar":"http://cl.ly/image/3e1h0H3H2s0P/concourse-logo.png","flow":"concourse-test","flow_token":"123","organization":"fairfaxmedia"},"params":{"event":"activity","message_body":"message body here","status_colour":"lime","status_value":"SUCCESS","title":"test title"}}`

	var input Input
	readInput([]byte(sample), &input)

	source := input.Source
	params := input.Params

	if strings.Compare(source.Author, "concourse") != 0 {
		t.Errorf("Author in Source was incorrect, expected: %s, got: %s", "concourse", source.Author)
	}

	if strings.Compare(params.StatusValue, "SUCCESS") != 0 {
		t.Errorf("Status in Params was incorrect, expected: %s, got: %s", "SUCCESS", params.StatusValue)
	}
}

func TestBuildRequestData(t *testing.T) {
	buildId := "1"
	buildName := "Test build"
	buildJobName := "Test build job"
	buildPipelineName := "Test pipeline"

	os.Setenv("BUILD_ID", buildId)
	os.Setenv("BUILD_NAME", buildName)
	os.Setenv("BUILD_JOB_NAME", buildJobName)
	os.Setenv("BUILD_PIPELINE_NAME", buildPipelineName)

	sample := `{"source":{"author":"concourse","avatar":"http://cl.ly/image/3e1h0H3H2s0P/concourse-logo.png", "flow_api":"http://api.flowdock.com", "flow_token":"123", "event": "message"},"params":{"event":"activity","message_body":"message body here","status_colour":"lime","status_value":"SUCCESS","title":"test title"}}`

	var input Input
	readInput([]byte(sample), &input)

	requestData := buildRequestData(&input)

	if strings.Compare(requestData["event"].(string), "activity") != 0 {
		t.Errorf("Request event tyoe was incorrect, expected: %s, got: %s", "activity", requestData["event"])
	}
	if strings.Compare(requestData["author"].(map[string]string)["name"], "concourse") != 0 {
		t.Errorf("Request author was incorrect, expected: %s, got: %s", "concourse", requestData["author"].(map[string]string)["name"])
	}
}

func TestExternalUrl(t *testing.T) {
	buildName := "1"
	buildJobName := "build"
	buildPipelineName := "pipeline"
	buildTeamName := "team"
	atcExternalUrl := "http://ci.url"

	os.Setenv("BUILD_NAME", buildName)
	os.Setenv("BUILD_TEAM_NAME", buildTeamName)
	os.Setenv("BUILD_JOB_NAME", buildJobName)
	os.Setenv("BUILD_PIPELINE_NAME", buildPipelineName)
	os.Setenv("ATC_EXTERNAL_URL", atcExternalUrl)

	expected := "http://ci.url/teams/team/pipelines/pipeline/jobs/build/builds/1"
	url := externalUrl()
	if strings.Compare(url, expected) != 0 {
		t.Errorf("Message body url was incorrect, expected: %s, got: %s", expected, url)
	}
}

func TestMessageTitle(t *testing.T) {
	buildName := "1"
	buildJobName := "build"
	buildPipelineName := "pipeline"

	os.Setenv("BUILD_NAME", buildName)
	os.Setenv("BUILD_JOB_NAME", buildJobName)
	os.Setenv("BUILD_PIPELINE_NAME", buildPipelineName)

	sample := `{"source":{},"params":{"message_title":"test title"}}`

	var input Input
	readInput([]byte(sample), &input)

	requestData := buildRequestData(&input)
	eventTitle := fmt.Sprintf("%s #%s", "test title", buildName)
	if strings.Compare(requestData["title"].(string), eventTitle) != 0 {
		t.Errorf("Request message title was incorrect, expected: %s, got: %s", eventTitle, requestData["title"])
	}
}

func TestEventMessageTitle(t *testing.T) {
	buildName := "1"
	buildJobName := "build"
	buildPipelineName := "pipeline"

	os.Setenv("BUILD_NAME", buildName)
	os.Setenv("BUILD_JOB_NAME", buildJobName)
	os.Setenv("BUILD_PIPELINE_NAME", buildPipelineName)

	sample := `{"source":{},"params":{"message_title":"", "status_value": "success"}}`

	var input Input
	readInput([]byte(sample), &input)

	requestData := buildRequestData(&input)
	eventTitle := fmt.Sprintf("%s | %s | %s [%s]", buildPipelineName, buildJobName, buildName, "success")
	if strings.Compare(requestData["title"].(string), eventTitle) != 0 {
		t.Errorf("Request message title was incorrect, expected: %s, got: %s", eventTitle, requestData["title"])
	}
}
