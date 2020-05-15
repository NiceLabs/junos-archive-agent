package main

type Configuration struct {
	FTP    *FTPConfigure    `json:"ftp,omitempty"`
	GitHub *GitHubConfigure `json:"github"`
}

type FTPConfigure struct {
	Port     int    `json:"port,omitempty"`
	Host     string `json:"host,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type GitHubConfigure struct {
	Owner    string `json:"owner"`
	Repo     string `json:"repo"`
	Token    string `json:"token"`
	Prefix   string `json:"prefix,omitempty"`
}
