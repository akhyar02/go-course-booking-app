{{template "base" .}}

{{define "css"}}
<style>
  .room-image {
    max-width: 50%;
  }
</style>
{{ end }}

{{define "content"}}

<div class="container">
  <div class="row">
    <div class="col">
      <img
        src="static/images/marjors-suite.png"
        class="img-fluid img-thumbnail mx-auto d-block room-image"
        alt="room image"
      />
    </div>
  </div>

  <div class="row">
    <div class="col">
      <h1 class="text-center mt-4">Major's Suite</h1>
      <p>
        Your home away form home, set on the majestic waters of the Atlantic
        Ocean, this will be a vacation to remember. Your home away form home,
        set on the majestic waters of the Atlantic Ocean, this will be a
        vacation to remember. Your home away form home, set on the majestic
        waters of the Atlantic Ocean, this will be a vacation to remember. Your
        home away form home, set on the majestic waters of the Atlantic Ocean,
        this will be a vacation to remember. Your home away form home, set on
        the majestic waters of the Atlantic Ocean, this will be a vacation to
        remember. Your home away form home, set on the majestic waters of the
        Atlantic Ocean, this will be a vacation to remember.
      </p>
    </div>
  </div>

  <div class="row">
    <div class="col text-center">
      <a
        id="check-availability-button"
        href="/search-availability?room=major-suites"
        class="btn btn-success"
        >Check Availability</a
      >
    </div>
  </div>
</div>

{{ end }}