{{template "base" .}}

{{define "content"}}

<div class="container">
  <div class="row">
    <div class="col">
      <p>This is the about page</p>

      <p>data from execute: {{index .StringMap "test"}}</p>
    </div>
  </div>
</div>

{{ end }}
