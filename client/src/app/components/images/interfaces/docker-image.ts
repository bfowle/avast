export interface DockerImage {
  id:          string;
  parentId:    string;
  repoTags:    string[];
  repoDigests: string[];
  created:     number;
  size:        string;
  virtualSize: string;
  labels:      any;  // make interface?
  children:    DockerImage[];
}

