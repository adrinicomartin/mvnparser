// MIT License
//
// Copyright (c) 2019 Aloïs Micard
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package mvnparser

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Represent a POM file
type MavenProject struct {
	XMLName              xml.Name             `xml:"project"`
	ModelVersion         string               `xml:"modelVersion"`
	Parent               Parent               `xml:"parent"`
	GroupId              string               `xml:"groupId"`
	ArtifactId           string               `xml:"artifactId"`
	Version              string               `xml:"version"`
	Packaging            string               `xml:"packaging"`
	Name                 string               `xml:"name"`
	Repositories         []Repository         `xml:"repositories>repository"`
	Properties           Properties           `xml:"properties"`
	DependencyManagement DependencyManagement `xml:"dependencyManagement"`
	Dependencies         []Dependency         `xml:"dependencies>dependency"`
	Profiles             []Profile            `xml:"profiles"`
	Build                Build                `xml:"build"`
	PluginRepositories   []PluginRepository   `xml:"pluginRepositories>pluginRepository"`
}

// Represent the parent of the project
type Parent struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
}

type Properties map[string]string

type propertyEntries struct {
	Entries []propertyEntry `xml:",any"`
}
type propertyEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (props *Properties) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var entries propertyEntries
	if err := d.DecodeElement(&entries, &start); err != nil {
		return err
	}
	for _, entry := range entries.Entries {
		if *props == nil {
			*props = make(Properties)
		}
		(*props)[entry.XMLName.Local] = entry.Value
	}
	return nil
}

type Property struct {
	Key   string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

// Represent a dependency of the project
type Dependency struct {
	XMLName    xml.Name    `xml:"dependency"`
	GroupId    string      `xml:"groupId"`
	ArtifactId string      `xml:"artifactId"`
	Version    string      `xml:"version"`
	Classifier string      `xml:"classifier"`
	Type       string      `xml:"type"`
	Scope      string      `xml:"scope"`
	Exclusions []Exclusion `xml:"exclusions>exclusion"`
}

// Represent an exclusion
type Exclusion struct {
	XMLName    xml.Name `xml:"exclusion"`
	GroupId    string   `xml:"groupId"`
	ArtifactId string   `xml:"artifactId"`
}

type DependencyManagement struct {
	Dependencies []Dependency `xml:"dependencies>dependency"`
}

// Represent a repository
type Repository struct {
	Id   string `xml:"id"`
	Name string `xml:"name"`
	Url  string `xml:"url"`
}

type Profile struct {
	Id    string `xml:"id"`
	Build Build  `xml:"build"`
}

type Build struct {
	// todo: final name ?
	Plugins []Plugin `xml:"plugins>plugin"`
}

type Plugin struct {
	XMLName    xml.Name `xml:"plugin"`
	GroupId    string   `xml:"groupId"`
	ArtifactId string   `xml:"artifactId"`
	Version    string   `xml:"version"`
	//todo something like: Configuration map[string]string `xml:"configuration"`
	// todo executions
}

// Represent a pluginRepository
type PluginRepository struct {
	Id   string `xml:"id"`
	Name string `xml:"name"`
	Url  string `xml:"url"`
}

//Parse a pom.xml file and return the MavenProject representing it.
func Parse(pomxmlPath string) (*MavenProject, error) {
	f, err := os.Open(pomxmlPath)
	if err != nil {
		return nil, fmt.Errorf("can't open file %s, %v", pomxmlPath, err)
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("can't read file %s, %v", pomxmlPath, err)
	}

	var project MavenProject
	if err := xml.Unmarshal(bytes, &project); err != nil {
		log.Fatalf("unable to unmarshal pom file. Reason: %s", err)
	}
	return &project, nil
}

//GetProperty with a particular key. Case insensitive.
func (mp *MavenProject) GetProperty(key string) (value string, exist bool) {
	for k, v := range mp.Properties {
		if strings.ToLower(k) == strings.ToLower(key) {
			return v, true
		}
	}
	return "", false
}
