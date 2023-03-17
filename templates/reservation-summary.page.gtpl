{{template "base" .}}

{{define "content"}}

<div class="container">
  <div class="row">
    <div class="col">
      <h1 class="mt-3">Reservation Summary</h1>

      <table class="table table-striped">
        <thead>
          <tr>
            <th scope="col">Room Type</th>
            <th scope="col">Start Date</th>
            <th scope="col">End Date</th>
            <th scope="col">First Name</th>
            <th scope="col">Last Name</th>
            <th scope="col">Email</th>
            <th scope="col">Phone</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>{{.Data.roomType}}</td>
            <td>{{.Data.startDate}}</td>
            <td>{{.Data.endDate}}</td>
            <td>{{.Data.firstName}}</td>
            <td>{{.Data.lastName}}</td>
            <td>{{.Data.email}}</td>
            <td>{{.Data.phone}}</td>
          </tr>
        </tbody>
    </div>
  </div>
</div>
{{ end }}
