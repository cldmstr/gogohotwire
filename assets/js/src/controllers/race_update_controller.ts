import {Controller} from '@stimulus/core';
import * as Turbo from "@hotwired/turbo";

export default class extends Controller {
  source: any;
  raceidTarget: any;

  static targets = [ "raceid" ];

  connect() {
    this.source = new EventSource("/races/" + this.raceidTarget.value + "/update");

    this.source.onmessage = (e: any) => {
      console.log(atob(e.data))
      Turbo.session.renderStreamMessage(atob(e.data));
    };

    this.source.onerror = (err: any) => {
      console.error("EventSource failed: ", err);
    };
  }

  disconnect() {
    this.source.close();
  }
}