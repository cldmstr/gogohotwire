import * as Turbo from "@hotwired/turbo";
import { Application } from "@stimulus/core";

import ToggleDriversController from './controllers/toggle_drivers_controller';
import RaceUpdateController from './controllers/race_update_controller';

const turbo = Turbo.start();
const application = Application.start();

application.register("toggle-drivers", ToggleDriversController);
application.register("race-update", RaceUpdateController);
