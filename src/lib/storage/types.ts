export type Maker = {
  id: string;
  name: string;
  logo: string;
  website: string;
};

export type MakerPie = {
  year: number;
  makerid: string;
  id: string;
  displayname: string;
  fresh: boolean;
  labels: string[];
  image_file: string;
  web_link: string;
  pack_count: number;
  pack_price_in_pence: number;
  validated: boolean;
};

export type MakerPieRanking = {
  year: number;
  makerid: string;
  pieid: string;
  userid: string;
  pastry: number;
  filling: number;
  topping: number;
  looks: number;
  value: number;
  notes?: string;
  last_updated?: string;
};

export enum StorageError {
  GenericError,
  NotFound,
  BadInput
}
