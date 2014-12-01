package main

import (
	"fmt"
	"log"

	"github.com/haklop/bazooka/commons/matrix"

	bazooka "github.com/haklop/bazooka/commons"
)

const (
	SourceFolder      = "/bazooka"
	OutputFolder      = "/bazooka-output"
	MetaFolder        = "/meta"
	BazookaConfigFile = ".bazooka.yml"
	TravisConfigFile  = ".travis.yml"
	Jdk               = "jdk"
)

func main() {
	file, err := bazooka.ResolveConfigFile(SourceFolder)
	if err != nil {
		log.Fatal(err)
	}

	conf := &ConfigJava{}
	err = bazooka.Parse(file, conf)
	if err != nil {
		log.Fatal(err)
	}
	buildTool, err := detectBuildTool(SourceFolder)
	if err != nil {
		log.Fatal(err)
	}

	mx := matrix.Matrix{
		Jdk: conf.JdkVersions,
	}

	if len(conf.JdkVersions) == 0 {
		mx[Jdk] = []string{"oraclejdk8"}

	}

	mx.IterAll(func(permutation map[string]string, counter string) {
		if err := manageJdkVersion(counter, conf, permutation[Jdk], buildTool); err != nil {
			log.Fatal(err)
		}
	}, nil)
}

func detectBuildTool(source string) (string, error) {
	exist, err := bazooka.FileExists(fmt.Sprintf("%s/build.gradle", source))
	if err != nil {
		return "", err
	}
	if exist {
		wrapperExist, err := bazooka.FileExists(fmt.Sprintf("%s/gradlew", source))
		if err != nil {
			return "", err
		}
		if wrapperExist {
			return "gradlew", nil
		} else {
			return "gradle", nil
		}
	}
	exist, err = bazooka.FileExists(fmt.Sprintf("%s/pom.xml", source))
	if err != nil {
		return "", err
	}
	if exist {
		return "maven", nil
	}
	return "ant", nil
}

func manageJdkVersion(counter string, conf *ConfigJava, version, buildTool string) error {
	conf.JdkVersions = []string{}
	setDefaultInstall(conf, buildTool)
	setDefaultScript(conf, buildTool)
	image, err := resolveJdkImage(version)
	conf.FromImage = image
	if err != nil {
		return err
	}
	err = bazooka.AppendToFile(fmt.Sprintf("%s/%s", MetaFolder, counter), fmt.Sprintf("%s: %s\n", Jdk, version), 0755)
	if err != nil {
		return err
	}

	return bazooka.Flush(conf, fmt.Sprintf("%s/.bazooka.%s.yml", OutputFolder, counter))
}

func setDefaultInstall(conf *ConfigJava, buildTool string) {
	if len(conf.Install) == 0 {
		instruction := switchDefaultInstall(buildTool)
		if len(instruction) != 0 {
			conf.Install = []string{instruction}
		}
	}
}

func switchDefaultInstall(buildTool string) string {
	switch buildTool {
	case "maven":
		return "mvn install -DskipTests=true"
	case "gradle":
		return "gradle assemble"
	case "gradlew":
		return "./gradlew assemble"
	default:
		//Do nothing by default for Ant
		return ""
	}
}

func setDefaultScript(conf *ConfigJava, buildTool string) {
	if len(conf.Script) == 0 {
		conf.Script = []string{switchDefaultScript(buildTool)}
	}
}

func switchDefaultScript(buildTool string) string {
	switch buildTool {
	case "maven":
		return "mvn test"
	case "gradle":
		return "gradle check"
	case "gradlew":
		return "./gradlew check"
	case "ant":
		return "ant test"
	default:
		return ""
	}
}

func resolveJdkImage(version string) (string, error) {
	//TODO extract this from db
	javaMap := map[string]string{
		"openjdk6":   "bazooka/runner-java:openjdk6",
		"openjdk7":   "bazooka/runner-java:openjdk7",
		"oraclejdk6": "bazooka/runner-java:oraclejdk6",
		"oraclejdk7": "bazooka/runner-java:oraclejdk7",
		"oraclejdk8": "bazooka/runner-java:oraclejdk8",
	}
	if val, ok := javaMap[version]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Unable to find Bazooka Docker Image for Java Runnner %s\n", version)
}
