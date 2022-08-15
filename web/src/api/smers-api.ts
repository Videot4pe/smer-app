import type { SmerDto } from "../models/smer";

import BaseCrud from "./client/base-crud";

export default BaseCrud<SmerDto, SmerDto>("smers");
