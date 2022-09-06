import { Controller } from '@stimulus/core';

export default class extends Controller {
    isVisible : boolean = false;

    connect() {
        this.element.addEventListener("click", this.onclick, false);
    }

    onclick() {
        let d = document.getElementById("drivers-column");
        let c = document.getElementById("content-column")
        if (d == null || c == null) {
            return
        }
        if (this.isVisible) {
            d.classList.add("invisible")
            c.classList.remove("col-8")
            c.classList.add("col-12")
        } else {
            d.classList.remove("invisible")
            c.classList.add("col-8")
            c.classList.remove("col-12")
        }
        this.isVisible = !this.isVisible;
    }
}