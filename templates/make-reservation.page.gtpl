{{template "base" .}}

{{define "content"}}

<div class="container">
  <div class="row">
    <div class="col">
      <h1 class="mt-3">Make Reservation</h1>

      <form method="post" action="/api/reservations" class="needs-validation" novalidate>
        <input type="hidden" name="csrf_token" value="{{.CSRFToken}}" />
        <div class="form-group mt-3">
          <label for="room_type">Room Type</label>
          <select
            class="form-control"
            id="room_type"
            name="room_type"
            required
          >
            <option value=null>Select Room Type</option>
            <option value="general quarters" {{if eq .Data.roomType "general_quarters" }}selected{{end}}>General Quarters</option>
            <option value="major suites" {{if eq .Data.roomType "major_suites" }}selected{{end}}>Major Suites </option>
          </select>
        </div>

        <div class="form-group" id="reservation-dates">
          <div class="row">
            <div class="col-md-6">
              <label for="start_date">Starting Date</label>
              <input
                required
                class="form-control"
                type="date"
                name="start_date"
                id="start_date"
                value="{{.Data.startDate}}"
              />
            </div>
            <div class="col-md-6">
              <label for="end_date">Ending Date</label>
              <input
                required
                class="form-control"
                type="date"
                name="end_date"
                id="end_date"
                value="{{.Data.endDate}}"
              />
            </div>
          </div>

        <div class="form-group mt-3">
          <label for="first_name">First Name:</label>
          <input
            class="form-control"
            id="first_name"
            autocomplete="off"
            type="text"
            name="first_name"
            value=""
            required
          />
        </div>

        <div class="form-group">
          <label for="last_name">Last Name:</label>
          <input
            class="form-control"
            id="last_name"
            autocomplete="off"
            type="text"
            name="last_name"
            value=""
            required
          />
        </div>

        <div class="form-group">
          <label for="email">Email:</label>
          <input
            class="form-control"
            id="email"
            autocomplete="off"
            type="email"
            name="email"
            value=""
            required
          />
        </div>

        <div class="form-group">
          <label for="phone">Phone:</label>
          <input
            class="form-control"
            id="phone"
            autocomplete="off"
            type="text"
            name="phone"
            value=""
            required
          />
        </div>

        <hr />
        <input type="submit" class="btn btn-primary" value="Make Reservation" />
      </form>
    </div>
  </div>
</div>
{{ end }}

{{define "js"}}
  <script>
    // Example starter JavaScript for disabling form submissions if there are invalid fields
    (function () {
      "use strict";

      window.addEventListener(
        "load",
        function () {
          // Fetch all the forms we want to apply custom Bootstrap validation styles to
          var forms = document.getElementsByClassName("needs-validation");

          // Loop over them and prevent submission
          var validation = Array.prototype.filter.call(forms, function (
            form
          ) {
            form.addEventListener(
              "submit",
              async function (event) {
                event.preventDefault();
                form.classList.add("was-validated");
                if (form.checkValidity() === false) {
                  event.stopPropagation();
                  return
                }
                const formData = new FormData(form);
                delete formData["csrf_token"]
                const response = await fetch("/api/reservations", {
                  method: "POST",
                  headers:{
                    'Content-Type': 'application/json',
                    'X-CSRF-Token': formData.get('csrf_token')
                  },
                  body: JSON.stringify(Object.fromEntries(formData))
                })
                const data = await response.json()
                const errors = Object.entries(data.errors)
                const responseErrors = document.querySelectorAll('.response-error')
                responseErrors.forEach((error) => {
                  error.classList.remove('is-invalid')
                  error.classList.remove('response-error')
                  error.parentNode.removeChild(error.parentNode.lastChild)
                })
                if (errors.length > 0) {
                  errors.forEach((field) => {
                    const inputField = document.getElementById(field[0])
                    inputField.classList.add('is-invalid')
                    inputField.classList.add('response-error')
                    const small = document.createElement('p')
                    small.classList.add('text-danger')

                    let errorMessage = field[1].join(', ')
                    if (field[1].length>1)
                      errorMessage = field[1].substring(0, errorMessage.length-2)
                    small.innerHTML = errorMessage
                    inputField.parentNode.appendChild(small)
                    inputField.onchange = () => {
                      inputField.classList.remove('is-invalid')
                      inputField.classList.remove('response-error')
                      inputField.parentNode.removeChild(inputField.parentNode.lastChild)
                    }
                  })
                  return
                }

                return window.location.href ="/reservation-summary"
              },
              false
            );
          });
        },
        false
      );
    })();
  </script>
{{end}}
