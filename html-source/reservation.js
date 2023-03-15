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
  return {
    error: false,
    message: "Success",
  }
  const response = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
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

  const responseData = await checkAvailability("/api/reservation", {
    start_date: startDateInput.value,
    end_date: endDateInput.value,
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
      window.location.href = "./make-reservation.html";
    },
  });
}

const reservationForm = document.getElementById("reservation_form");
const startDateInput = document.getElementById("start_date");
const endDateInput = document.getElementById("end_date");
reservationForm.onsubmit = handleSubmit