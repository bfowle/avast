export interface DockerContainer {
  Id:               string;
  Names:            string[];
  Image:            string;
  ImageID:          string;
  Command:          string;
  Created:          number;
  Ports:            number[];
  Labels:           any;   // make interface?
  Status:           string;
  HostConfig:       any;   // make interface?
  NetworkSettings:  any;   // make interface?
}
