# Product Requirement Document (PRD) #001 

## Document Metadata

- Title : Initial Setup of Validra Project
- Author : Arif Setyawan
- Version : 0.0.1

## Overview

### Purpose

This PRD solely purposed to define base structure and functionality of services that capable to provide RBAC, ABAC and ReBAC functionalities. Focus on building fundamental like architecture, data structure and infrastructure. Initial setup on code structure and infrastrcture is priority in this PRD

### Problem Statement

Build from scratch. Initial setup on code structure and infrastrcture is priority in this PRD

### Goals/Objective

- Create and generate mature model, repository and interface API that would make validra works. Expecting API is functioning.
- Able to CRUD the `resources`, `role`, `user`, `user set` and `resource set`
- User set and resource set have condition group where it would accept chained of conditions.
- Functionality to /check-permission to check if the user with given criteria allowed. response as minimum is boolean. would be good if have clearer context for allow or rejecting the access.

### Vision

To provide service that could help segregation and independent to smartly enough maintain and manage RBAC, ABAC and ReBAC

## Product Description

### Summary

Validra engine is the open source software that do RBAC, ABAC and ReBAC. This software is written in golang and use sqlite as database. Initial setup on code structure and infrastrcture is priority in this PRD

### Target Audience 

- API users.

### Use Cases

- Api User Create Read Update and Delete the Resource
- Api User Create Read Update and Delete the Resource Set
- Api User Create Read Update and Delete the Role
- Api User Create Read Update and Delete the User
- Api User Create Read Update and Delete the User Set
- Api User to check access 

## Functional Requirements

- base API functionalities

## Technical Consideration

### Technologies Stack

- latest golang
- latest sqlite
- golang echo

### Database Schema
  - Resource
  - Role
  - User
  - User Set
  - Resource Set



