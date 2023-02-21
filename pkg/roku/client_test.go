package roku

import (
	"fmt"
	"testing"

	"github.com/picatz/roku"
)

func TestSomething(t *testing.T) {
	fmt.Println("")
	// client, err := NewClient(&ClientOptions{
	// 	Logger: log.Default(),
	// 	Config: &Config{
	// 		DeviceName: "BKYLN 57 Fifth TV",
	// 		WaitTime: 10 * time.Second,
	// 	},
	// })
	// require.NoError(t, err)

	// client.StartApp(&StartAppOptions{
	// 	AppName: "Live TV",
	// 	Device:  roku.NewEndpoint("http://192.168.1.18:8060/"),
	// 	Force:   true,
	// })

	// results, err := client.Search()
	// require.NoError(t, err)
	// fmt.Println(results)

	TV := roku.NewEndpoint("http://192.168.1.18:8060/")
	// info, err := TV.DeviceInfo()
	// require.NoError(t, err)
	// fmt.Println(info.FriendlyDeviceName)

	err := TV.LaunchApp("08532417380d48cf2244a1831ff94f26", nil)
	fmt.Println(err)
	// err := TV.PlayVideo("http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/SubaruOutbackOnStreetAndDirt.mp4")
	// require.NoError(t, err)

	// apps, err := TV.Apps()
	// require.NoError(t, err)
	// for _, app := range apps {
	// 	fmt.Println(app.Name, app.ID)
	// }
}
