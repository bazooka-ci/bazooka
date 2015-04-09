package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	bazooka "github.com/bazooka-ci/bazooka/commons"
	bzklog "github.com/bazooka-ci/bazooka/commons/logs"
)

const (
	SourceFolder = "/bazooka"
	OutputFolder = "/bazooka-output"
	MetaFolder   = "/meta"
	Jdk          = "jdk"
)

func init() {
	log.SetFormatter(&bzklog.BzkFormatter{})
}

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

	versions := conf.JdkVersions
	images := conf.Base.Image

	if len(versions) == 0 && len(images) == 0 {
		versions = []string{"oraclejdk8"}

	}

	for i, version := range versions {
		if err := manageJdkVersion(fmt.Sprintf("0%d", i), conf, version, "", buildTool); err != nil {
			log.Fatal(err)
		}
	}

	for i, image := range images {
		if err := manageJdkVersion(fmt.Sprintf("1%d", i), conf, "", image, buildTool); err != nil {
			log.Fatal(err)
		}
	}
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

func manageJdkVersion(counter string, conf *ConfigJava, version, image, buildTool string) error {
	conf.JdkVersions = nil
	conf.Base.Image = nil

	setDefaultInstall(conf, buildTool)
	setDefaultScript(conf, buildTool)

	meta := map[string]string{}
	if len(version) > 0 {
		var err error
		image, err = resolveJdkImage(version)
		if err != nil {
			return err
		}
		meta[Jdk] = version
	} else {
		meta["image"] = image
	}
	conf.Base.FromImage = image

	if err := bazooka.Flush(meta, fmt.Sprintf("%s/%s", MetaFolder, counter)); err != nil {
		return err
	}

	return bazooka.Flush(conf, fmt.Sprintf("%s/.bazooka.%s.yml", OutputFolder, counter))
}

func setDefaultInstall(conf *ConfigJava, buildTool string) {
	if len(conf.Base.Install) == 0 {
		instruction := switchDefaultInstall(buildTool)
		if len(instruction) != 0 {
			conf.Base.Install = []string{instruction}
		}
	}
}

func switchDefaultInstall(buildTool string) string {
	switch buildTool {
	case "maven":
		return "mvn install -DskipTests=true --batch-mode"
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
	if len(conf.Base.Script) == 0 {
		conf.Base.Script = []string{switchDefaultScript(buildTool)}
	}
}

func switchDefaultScript(buildTool string) string {
	switch buildTool {
	case "maven":
		return "mvn test --batch-mode"
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
		"openjdk8":   "bazooka/runner-java:openjdk8",
		"oraclejdk6": "bazooka/runner-java:oraclejdk6",
		"oraclejdk7": "bazooka/runner-java:oraclejdk7",
		"oraclejdk8": "bazooka/runner-java:oraclejdk8",
	}
	if val, ok := javaMap[version]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Unable to find Bazooka Docker Image for Java Runnner %s\n", version)
}
