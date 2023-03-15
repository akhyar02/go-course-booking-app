{{template "base" .}}

{{define "content"}}

<div class="container">
  <div class="row">
    <div class="col">
      <h1 class="mt-3">Make Reservation</h1>

      <form method="post" action="" class="needs-validation" novalidate>
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
              function (event) {
                if (form.checkValidity() === false) {
                  event.preventDefault();
                  event.stopPropagation();
                }
                form.classList.add("was-validated");
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
