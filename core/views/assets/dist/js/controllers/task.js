import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
    connect() {
        console.log(this.element)
    }

    toggle(e) {
        const id = e.target.dataset.id

        fetch(`/tasks/${id}/toggle`, {
            method: 'POST', // *GET, POST, PUT, DELETE, etc.
            mode: 'cors', // no-cors, *cors, same-origin
            cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
            credentials: 'same-origin', // include, *same-origin, omit
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ completed: e.target.checked }) // body data type must match "Content-Type" header
        })
            .then(response => response.json())
            .then(data => {
                alert(data.message)
            })
    }
}

