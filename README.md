# Instagrab

Little utility for download photo and video from social media `Instagram.com` in max direction.

## Setup

### Dependencies

```sh
$ sudo apt-get update && sudo apt-get upgrade
$ sudo apt-get install golang-go
$ sudo apt-get install git
```

### Build

```sh
$ git clone https://gitlab.com/N1ghtRaven/instagrab.git
$ go build instagrab.go
```

## Usage
### Print direct-link on console

```sh
$ ./instagrab -only-url (-url|-showcode) [URL_OR_SHOWCODE]  
```
### Download

```sh
$ ./instagrab (-url|-showcode) [URL_OR_SHOWCODE]  
```