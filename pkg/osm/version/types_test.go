package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"testing"

	tassert "github.com/stretchr/testify/assert"
)

func TestEnvoyConfigParser(t *testing.T) {
	assert := tassert.New(t)
	actual := getReleases()

	assert.Equal(actual, []string{"v0.10", "v0.11", "v0.5", "v0.6", "v0.7", "v0.8", "v0.9"})

	for _, release := range actual {
		controllerVersion := ControllerVersion(release)
		{
			_, exists := SupportedIngress[controllerVersion]
			assert.Truef(exists, "IngressVersion does not contain info on OSM release %s", release)
		}

		{
			_, exists := SupportedTrafficTarget[controllerVersion]
			assert.Truef(exists, "SupportedTrafficTarget does not contain info on OSM release %s", release)
		}

		{
			_, exists := SupportedTrafficTargetRouteKinds[controllerVersion]
			assert.Truef(exists, "SupportedTrafficTargetRouteKinds does not contain info on OSM release %s", release)
		}

		{
			_, exists := SupportedTrafficSplit[controllerVersion]
			assert.Truef(exists, "SupportedTrafficSplit does not contain info on OSM release %s", release)
		}

		{
			_, exists := SupportedHTTPRouteVersion[controllerVersion]
			assert.Truef(exists, "SupportedHTTPRouteVersion does not contain info on OSM release %s", release)
		}

		{
			_, exists := SupportedAnnotations[controllerVersion]
			assert.Truef(exists, "SupportedAnnotations does not contain info on OSM release %s", release)
		}

		{
			_, exists := EnvoyAdminPort[controllerVersion]
			assert.Truef(exists, "EnvoyAdminPort does not contain info on OSM release %s", release)
		}

		{
			_, exists := OutboundListenerNames[controllerVersion]
			assert.Truef(exists, "OutboundListenerNames does not contain info on OSM release %s", release)
		}

		{
			_, exists := InboundListenerNames[controllerVersion]
			assert.Truef(exists, "InboundListenerNames does not contain info on OSM release %s", release)
		}
	}
}

const releasesURL = "https://api.github.com/repos/openservicemesh/osm/releases"

func getReleases() []string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", releasesURL, nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal().Err(err)
		}
	}()

	var res []interface{}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Fatal().Err(err)
	}

	ignore := map[string]interface{}{
		"v0.1": nil,
		"v0.2": nil,
		"v0.3": nil,
		"v0.4": nil,
	}
	releases := make(map[string]interface{})

	for _, releaseJSON := range res {
		rel := releaseJSON.(map[string]interface{})
		tagNameIface, ok := rel["tag_name"]
		if !ok {
			continue
		}
		tagName, ok := tagNameIface.(string)
		if !ok || strings.Contains(tagName, "-rc") {
			continue
		}

		majorMinorChunks := strings.Split(tagName, ".")
		if len(majorMinorChunks) < 2 {
			continue
		}
		release := fmt.Sprintf("%s.%s", majorMinorChunks[0], majorMinorChunks[1])
		if _, shouldIgnore := ignore[release]; shouldIgnore {
			continue
		}
		releases[release] = nil
	}

	var releasesSlice []string
	for rel := range releases {
		releasesSlice = append(releasesSlice, rel)
	}

	sort.Strings(releasesSlice)

	return releasesSlice
}
