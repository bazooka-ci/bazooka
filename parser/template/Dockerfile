FROM {{.Generator.Config.FromImage}}

COPY source {{.BzkBuildDir}}

COPY work/{{$.Generator.Index}}/bazooka_*.sh {{$.BzkBuildDir}}/
RUN chmod +x {{$.BzkBuildDir}}/bazooka_*.sh

{{range .Generator.Config.Env}}
ENV {{.Name}} {{.Value}}
{{end}}

WORKDIR {{.BzkBuildDir}}

CMD  {{.BzkBuildDir}}/bazooka_run.sh
