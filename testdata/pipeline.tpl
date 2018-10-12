{{- set . "resource_types" (ds "pipeline-resources.yaml").resource_types -}}
{{- set . "resources" (ds "pipeline-resources.yaml").resources -}}
{{- .|yaml}}
