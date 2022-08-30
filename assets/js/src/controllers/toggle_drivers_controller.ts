import { Controller } from '@stimulus/core';
import { Collapse } from 'bootstrap';

export default class extends Controller {
    isVisible : boolean = false;

    connect() {
        this.element.addEventListener("click", this.onclick, false);
        console.log("Connected Toggle Drivers");
        console.log(this.element);
    }

    onclick() {
        let d = document.getElementById("drivers-column");
        let c = document.getElementById("content-column")
        if (d == null || c == null) {
            return
        }
        if (this.isVisible) {
            d.classList.add("invisible")
            c.classList.remove("col-9")
            c.classList.add("col-12")
        } else {
            d.classList.remove("invisible")
            c.classList.add("col-9")
            c.classList.remove("col-12")
        }
        this.isVisible = !this.isVisible;
    }
}