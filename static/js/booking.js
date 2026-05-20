// Shared availability check and Book Now flow (requires Swal + attention from base layout)

function checkAvailability(form, roomId, roomName, csrfToken) {
    const formData = new FormData(form);
    formData.append('csrf_token', csrfToken);
    if (roomId) {
        formData.append('room_id', roomId);
    }
    if (roomName) {
        formData.append('room_name', roomName);
    }

    return fetch('/search-availability-json', {
        method: 'post',
        body: formData,
    }).then(response => response.json());
}

function showAvailabilityResult(data, csrfToken) {
    if (!data.ok) {
        attention.error({
            title: 'Not available',
            msg: data.message || 'No availability for those dates.',
        });
        return;
    }

    const roomLine = data.room_name
        ? `<p><strong>Room:</strong> ${data.room_name}</p>`
        : '';

    Swal.fire({
        icon: 'success',
        title: 'Dates confirmed',
        html: `
            <p><strong>Arrival:</strong> ${data.start_date}</p>
            <p><strong>Departure:</strong> ${data.end_date}</p>
            ${roomLine}
            <button type="button" id="book-now-btn" class="btn btn-primary btn-block mt-3">Book Now</button>
        `,
        showConfirmButton: false,
        showCloseButton: true,
        didOpen: () => {
            document.getElementById('book-now-btn').addEventListener('click', () => {
                Swal.close();
                promptEmailAndBook(data, csrfToken);
            });
        },
    });
}

function promptEmailAndBook(availability, csrfToken) {
    Swal.fire({
        title: 'Book Now',
        html: `
            <p class="text-left mb-2"><strong>Arrival:</strong> ${availability.start_date}</p>
            <p class="text-left mb-3"><strong>Departure:</strong> ${availability.end_date}</p>
        `,
        input: 'email',
        inputLabel: 'Your email address',
        inputPlaceholder: 'you@example.com',
        showCancelButton: true,
        confirmButtonText: 'Send confirmation',
        preConfirm: (email) => {
            if (!email) {
                Swal.showValidationMessage('Please enter your email address');
                return false;
            }
            const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!re.test(email)) {
                Swal.showValidationMessage('Please enter a valid email address');
                return false;
            }
            return email;
        },
    }).then((result) => {
        if (!result.isConfirmed) {
            return;
        }

        const formData = new FormData();
        formData.append('csrf_token', csrfToken);
        formData.append('email', result.value);
        formData.append('start', availability.start_date);
        formData.append('end', availability.end_date);
        if (availability.room_id) {
            formData.append('room_id', availability.room_id);
        }
        if (availability.room_name) {
            formData.append('room_name', availability.room_name);
        }

        fetch('/book-now', {
            method: 'post',
            body: formData,
        })
            .then(response => response.json())
            .then(booking => {
                if (booking.ok) {
                    Swal.fire({
                        icon: 'success',
                        title: 'Booking confirmed',
                        text: booking.message,
                        confirmButtonText: 'View details',
                    }).then(() => {
                        window.location.href = '/reservation-summary';
                    });
                } else {
                    attention.error({
                        title: 'Booking failed',
                        msg: booking.message || 'Could not complete your booking.',
                    });
                }
            })
            .catch(() => {
                attention.error({
                    title: 'Booking failed',
                    msg: 'Something went wrong. Please try again.',
                });
            });
    });
}

function bindAvailabilityCheck(buttonId, roomId, roomName, csrfToken) {
    document.getElementById(buttonId).addEventListener('click', function () {
        const html = `
        <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
            <div class="form-row">
                <div class="col">
                    <div class="form-row" id="reservation-dates-modal">
                        <div class="col">
                            <input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival">
                        </div>
                        <div class="col">
                            <input disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure">
                        </div>
                    </div>
                </div>
            </div>
        </form>
        `;

        attention.custom({
            title: 'Choose your dates',
            msg: html,
            willOpen: () => {
                const elem = document.getElementById('reservation-dates-modal');
                new DateRangePicker(elem, {
                    format: 'yyyy-mm-dd',
                    showOnFocus: true,
                });
            },
            didOpen: () => {
                document.getElementById('start').removeAttribute('disabled');
                document.getElementById('end').removeAttribute('disabled');
            },
            callback: function (result) {
                if (!result) {
                    return;
                }
                const form = document.getElementById('check-availability-form');
                checkAvailability(form, roomId, roomName, csrfToken)
                    .then(data => showAvailabilityResult(data, csrfToken));
            },
        });
    });
}
