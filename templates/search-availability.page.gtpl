{{template "base" .}}

{{define "content"}}

<div class="container">
  <div class="row">
    <div class="col">
      <h1 class="mt-3">Search for Availability</h1>

      <form
        novalidate
        class="needs-validation"
        id="reservation_form"
      >
        <div class="form-group">
          <label for="room_type">Room Type</label>
          <select
            class="form-control"
            id="room_type"
            name="room_type"
            required
          >
            <option value=null>Select Room Type</option>
            <option value="generalQuarters" {{
              if eq .Data.roomType "general-quarters" }}selected{{end}}>General Quarters</option>
            <option value="majorSuites" {{
              if eq .Data.roomType "major-suites" }}selected{{end}}>Major Suites </option>
          </select>
        <input type="hidden" name="csrf_token" value="{{.CSRFToken}}" />
        <div class="row">
          <div class="col">
            <div class="row" id="reservation-dates">
              <div class="col-md-6">
                <label for="start_date">Starting Date</label>
                <input
                  required
                  class="form-control"
                  type="date"
                  name="start_date"
                  id="start_date"
                />
                <small id="startDateHelp" class="form-text text-muted"
                  >Enter your starting date in YYYY-MM-DD format</small
                >
              </div>
              <div class="col-md-6">
                <label for="end_date">Ending Date</label>
                <input
                  required
                  class="form-control"
                  type="date"
                  name="end_date"
                  id="end_date"
                />
                <small id="endDateHelp" class="form-text text-muted"
                  >Enter your ending date in YYYY-MM-DD format</small
                >
              </div>
            </div>
          </div>
        </div>

        <hr />

        <button type="submit" class="btn btn-primary">
          Search Availability
        </button>

        <div class="invalid-feedback">
          Please provide a valid starting and ending date.
        </div>
      </form>
    </div>
  </div>
</div>

{{ end }}

{{define "js"}}

<script src="https://cdn.jsdelivr.net/npm/sweetalert2@11"></script>
<script>
  function validateInput(input) {
  if (!input.value) {
    input.classList.add("is-invalid");
    return false;
  }
  input.classList.remove("is-invalid");
  return true;
}

function validateDate(startDate, endDate) {
  if (startDate.value > endDate.value) {
    startDate.classList.add("is-invalid");
    endDate.classList.add("is-invalid");
    document.querySelector(".invalid-feedback").style.display = "block";
    return false;
  }
  startDate.classList.remove("is-invalid");
  endDate.classList.remove("is-invalid");
  document.querySelector(".invalid-feedback").style.display = "none";
  return true;
}

async function checkAvailability(url = "", data = {}) {
  // return {
  //   error: false,
  //   message: "Success",
  // }
  const query = `?startDate=${data.startDate}&endDate=${data.endDate}`
  const response = await fetch(url+query, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    }
  });
  if (!response.ok) {
    return {
      error: true,
      message: `An error has occured: ${response.status}`,
    };
  }
  return response.json();
}

async function handleSubmit(e){
  e.preventDefault();
  if (roomTypeInput.value === "null") {
    roomTypeInput.classList.add("is-invalid");
    document.querySelector(".invalid-feedback").style.display = "none";
    return;
  } else {
    roomTypeInput.classList.remove("is-invalid");
  }
  if (
    !validateInput(startDateInput) ||
    !validateInput(endDateInput) ||
    !validateDate(startDateInput, endDateInput)
  ) {
    setTimeout(() => {
      startDateInput.classList.remove("is-invalid");
      endDateInput.classList.remove("is-invalid");
      document.querySelector(".invalid-feedback").style.display = "none";
    }, 5000);
    return;
  }

  const responseData = await checkAvailability("/api/reservations", {
    startDate: startDateInput.value,
    endDate: endDateInput.value,
  });
  if (responseData.error) {
    Swal.fire({
      title: "Error!",
      text: responseData.message,
      icon: "error",
      confirmButtonText: "OK",
    });
    return;
  }
  Swal.fire({
    title: "Room is available!",
    icon: "success",
    confirmButtonText: "Make reservation",
    preConfirm: () => {
      window.location.href = "/make-reservation";
    },
  });
}

const reservationForm = document.getElementById("reservation_form");
const startDateInput = document.getElementById("start_date");
const endDateInput = document.getElementById("end_date");
const roomTypeInput = document.getElementById("room_type");
reservationForm.onsubmit = handleSubmit
</script>

{{ end }}
