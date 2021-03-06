package enumerate

import (
	"encoding/json"
	"fmt"
	"regexp"

	. "github.com/blackcrw/wprecon/cli/config"
	"github.com/blackcrw/wprecon/pkg/printer"
	"github.com/blackcrw/wprecon/pkg/text"
	"github.com/blackcrw/wprecon/tools/wordpress/extensions"
)

type uJSON []struct {
	Name string `json:"name"`
}

// UsersEnumeratePassive :: Enumerate using feed
func UsersEnumeratePassive() []string {
	response := extensions.SimpleRequest(InfosWprecon.Target, "feed/")

	rex := regexp.MustCompile("<dc:creator><!\\[CDATA\\[(.+?)\\]\\]></dc:creator>")
	submatch := rex.FindAllSubmatch([]byte(response.Raw), -1)

	for _, value := range submatch {
		valueString := fmt.Sprintf("%s", value[1])

		if _, has := text.ContainsSliceString(InfosWprecon.OtherInformationsSlice["target.http.users"], valueString); !has {
			InfosWprecon.OtherInformationsSlice["target.http.users"] = append(InfosWprecon.OtherInformationsSlice["target.http.users"], valueString)
		}
	}

	if len(InfosWprecon.OtherInformationsSlice["target.http.users"]) > 0 {
		InfosWprecon.OtherInformationsString["target.http.users.method"] = "Feed"
	}

	return InfosWprecon.OtherInformationsSlice["target.http.users"]
}

// UsersEnumerateAgressive ::
func UsersEnumerateAgressive() []string {
	var ujson uJSON
	done := false

	// Enumerate using Yoast SEO
	func() {
		if done == false {
			response := extensions.SimpleRequest(InfosWprecon.Target, "author-sitemap.xml")

			rex := regexp.MustCompile("<loc>.*?/author/(.*?)/</loc>")

			submatch := rex.FindAllSubmatch([]byte(response.Raw), -1)

			for _, value := range submatch {
				valueString := fmt.Sprintf("%s", value[1])

				if _, has := text.ContainsSliceString(InfosWprecon.OtherInformationsSlice["target.http.users"], valueString); !has {
					InfosWprecon.OtherInformationsSlice["target.http.users"] = append(InfosWprecon.OtherInformationsSlice["target.http.users"], valueString)
				}
			}

			if len(InfosWprecon.OtherInformationsSlice["target.http.users"]) > 0 {
				InfosWprecon.OtherInformationsString["target.http.users.method"] = "YoastSEO"
				done = true
			}
		}
	}()

	// Enumerate using route
	func() {
		if done == false {
			response := extensions.SimpleRequest(InfosWprecon.Target, "?rest_route=/wp/v2/users")

			if response.StatusCode == 200 && response.Raw != "" {
				json.NewDecoder(response.RawIo).Decode(&ujson)

				for _, value := range ujson {
					if _, has := text.ContainsSliceString(InfosWprecon.OtherInformationsSlice["target.http.users"], value.Name); !has {
						InfosWprecon.OtherInformationsSlice["target.http.users"] = append(InfosWprecon.OtherInformationsSlice["target.http.users"], value.Name)
					}
				}

				if len(InfosWprecon.OtherInformationsSlice["target.http.users"]) > 0 {
					InfosWprecon.OtherInformationsString["target.http.users.method"] = "Route"
					done = true
				}
			} else if response.StatusCode == 401 && response.Raw != "" && InfosWprecon.Verbose {
				printer.Danger("Status code 401, I don't think I'm allowed to list users. Target Url:", response.URL.Full, "— Target source code:", response.Raw)
			}
		}
	}()

	// Enumerate using json file
	func() {
		if done == false {
			response := extensions.SimpleRequest(InfosWprecon.Target, "wp-json/wp/v2/users")

			if response.StatusCode == 200 && response.Raw != "" {
				json.NewDecoder(response.RawIo).Decode(&ujson)

				for _, value := range ujson {
					if _, has := text.ContainsSliceString(InfosWprecon.OtherInformationsSlice["target.http.users"], value.Name); !has {
						InfosWprecon.OtherInformationsSlice["target.http.users"] = append(InfosWprecon.OtherInformationsSlice["target.http.users"], value.Name)
					}
				}

				if len(InfosWprecon.OtherInformationsSlice["target.http.users"]) > 0 {
					InfosWprecon.OtherInformationsString["target.http.users.method"] = "Json"
					done = true
				}
			} else if response.StatusCode == 401 && response.Raw != "" && InfosWprecon.Verbose {
				printer.Danger("Status code 401, I don't think I'm allowed to list users. Target Url:", response.URL.Full, "— Target source code:", response.Raw)
			}
		}
	}()

	return InfosWprecon.OtherInformationsSlice["target.http.users"]
}
