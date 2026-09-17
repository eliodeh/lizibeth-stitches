// ── SET MIN DATE TO TODAY ──
document.addEventListener('DOMContentLoaded', function() {
    const dateInput = document.getElementById('booking-date')
    if (dateInput) {
        const today = new Date().toISOString().split('T')[0]
        dateInput.min = today

        // load slots when date or service changes
        dateInput.addEventListener('change', loadSlots)
    }

    // listen for service radio change
    const serviceRadios = document.querySelectorAll('input[name="service"]')
    serviceRadios.forEach(radio => {
        radio.addEventListener('change', loadSlots)
    })
})

// ── LOAD AVAILABLE SLOTS ──
function loadSlots() {
    const dateInput   = document.getElementById('booking-date')
    const timeSlot    = document.getElementById('time-slot')
    const serviceRadio = document.querySelector('input[name="service"]:checked')

    if (!dateInput || !timeSlot) return

    const date    = dateInput.value
    const service = serviceRadio ? serviceRadio.value : ''

    if (!date || !service) {
        timeSlot.innerHTML = '<option value="">Select date and service first</option>'
        return
    }

    // check if Sunday
    const selectedDate = new Date(date)
    if (selectedDate.getDay() === 0) {
        timeSlot.innerHTML = '<option value="">Sundays are not available</option>'
        return
    }

    // fetch available slots from server
    timeSlot.innerHTML = '<option value="">Loading slots...</option>'

    fetch(`/booking/slots?date=${date}&service=${service}`)
        .then(response => response.text())
        .then(html => {
            timeSlot.innerHTML = html
        })
        .catch(err => {
            console.error('Error loading slots:', err)
            timeSlot.innerHTML = '<option value="">Error loading slots</option>'
        })
}

// ── FOR RESCHEDULE PAGE ──
// load slots with service from hidden booking data
const rescheduleService = document.querySelector('input[name="service_type"]')
if (rescheduleService) {
    const dateInput = document.getElementById('booking-date')
    if (dateInput) {
        dateInput.addEventListener('change', function() {
            const timeSlot = document.getElementById('time-slot')
            const date = this.value
            const service = rescheduleService.value

            if (!date) return

            fetch(`/booking/slots?date=${date}&service=${service}`)
                .then(response => response.text())
                .then(html => {
                    timeSlot.innerHTML = html
                })
        })
    }
}