package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
	version "github.com/hashicorp/go-version"
)

func currentVersion(asin_number string) *version.Version {
	// Request the HTML page.

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://www.amazon.co.jp/dp/"+asin_number, nil)

	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36")

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// bodyBytes, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// bodyString := string(bodyBytes)
	// log.Println(bodyString)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	versionStr := doc.Find("#masTechnicalDetails-btf > div:nth-child(4) > span:nth-child(2)").Text()
	versionStr = strings.TrimSpace(versionStr)

	if versionStr == "" {
		log.Fatal("No version has been found")
	}

	v1, err := version.NewVersion(versionStr)
	if err != nil {
		log.Fatal("Version is not parsable")
	}

	return v1
}

func main() {

	asin_number := strings.TrimSpace(os.Getenv("asin_number"))

	if asin_number == "" {
		fmt.Println(" ASIN Number is invalid. Exiting...")
		os.Exit(1)
	}

	fmt.Println("'ASIN' number is :", asin_number)

	currentVersion := currentVersion(asin_number)
	log.Printf("Version is %x \n", currentVersion)
	segments := currentVersion.Segments()
	updatedVersion := fmt.Sprintf("%d.%d.%d", segments[0], segments[1], segments[2]+1)
	log.Printf("Updated version is  %s \n", updatedVersion)

	fmt.Println("Amazon client id is:", os.Getenv("amazon_client_id"))
	fmt.Println("Amazon client secret is:", os.Getenv("amazon_client_secret"))

	cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", "AMAZON_RELEASE_VERSION", "--value", updatedVersion).CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
		os.Exit(1)
	}

	os.Exit(0)
}
