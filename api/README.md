# API Documentation

## Introduction

This document will introduce the different APIs exposed by this service and how users can interact with them.

### APOD Images

This returns a batch of APOD objects. You must specify the `count` and `page` values.

- `count`: `int`
  - How many APOD objects you want to return in a single page
- `page`: `int`
  - The page you want to return

Example Usage. This will return the first 30 images from today's date

```
GET Request: localhost:8081/api/v1/apod/batch/?count=30&page=0
```
