# RODAMAP.MVP.UI.md

> **Project:** MES Lite - Web UI Roadmap
>
> This document tracks the server-rendered Web UI needed for the Excel/paper-replacement MVP.
>
> `ROADMAP.md` remains API/backend focused. This file covers browser screens, forms, navigation and role-specific UI flows built on top of the existing API/application server.

## UI Goal

Build a simple server-rendered web interface that lets a small manufacturing company replace current Excel/paper production tracking.

The UI should be part of the Go application server. It should not introduce a separate frontend application for the MVP.

## Scope

The UI MVP should expose only the workflows needed to replace current manual tracking:

- login/logout
- worker production registration
- manager review of entered production records
- manager access to manual-summary reports
- basic master-data screens needed to support registration
- CSV import screen for historical manual data

Anything beyond replacing current manual tracking belongs to MVP-V2 or later.

## UX Modes

The same application should provide different UX modes rather than separate applications.

### Worker Mode

Worker mode should be optimized for fast, low-error data entry:

- large controls
- minimal navigation
- clear current form state
- minimal typing
- good tablet/kiosk behavior
- simple recovery from validation errors

### Manager/Admin Mode

Manager/admin mode can use denser screens:

- tables
- filters
- detail pages
- report views
- administrative forms

## UI Milestones

### UI Milestone 1 - Web Foundation

Goal: add the minimal server-rendered UI foundation inside the Go application.

Features:

- HTML layout
- route group for UI pages
- session-aware navigation
- login/logout pages
- shared error display pattern
- basic access control for UI routes

Definition of Done:

- authenticated users can access role-appropriate pages
- unauthenticated users are redirected or rejected consistently
- UI routes do not duplicate business logic from API/domain services

### UI Milestone 2 - Worker Production Registration

Goal: let workers register production without Excel or paper.

Features:

- production registration form
- employee/product/order selection using existing data
- workstation text field
- quantity entry
- timestamp handling
- comment/problem field if currently tracked manually
- validation error display
- success confirmation

Definition of Done:

- worker can enter production records from the browser
- submitted records use existing production-registration backend logic
- duplicate submissions are handled through the API idempotency/request-ID behavior

### UI Milestone 3 - Manager Production Review

Goal: let managers and leaders review entered production records.

Features:

- production-entry list
- simple filters matching existing API capabilities
- production-entry detail view if needed for correction/review
- link to correction workflow when permitted

Definition of Done:

- managers/leaders can review production records without querying the database or exporting raw data manually
- workers cannot access manager review screens

### UI Milestone 4 - Corrections UI

Goal: let permitted users correct production-entry mistakes without silently editing history.

Features:

- correction form
- reason field
- display of original entry
- display of correction history where available

Definition of Done:

- leaders/managers/admins can submit corrections through the browser
- normal workers cannot correct historical records
- corrections use the API correction model and preserve audit history

### UI Milestone 5 - Reports UI

Goal: expose manual-summary replacement reports in the browser.

Features:

- daily production report page
- employee productivity report page where manually tracked
- product statistics report page where manually tracked
- date-range filters

Definition of Done:

- managers/leaders can view manual-summary reports without maintaining Excel reports manually
- report pages use existing reporting endpoints/read models

### UI Milestone 6 - CSV Import UI

Goal: let authorized users import historical manual data.

Features:

- CSV upload form
- validation summary display
- import result display
- partial failure/error display

Definition of Done:

- authorized users can upload historical CSV data from the browser
- import errors are understandable without reading server logs

## Future UI Scope

The following are not required for the UI MVP unless they directly replace current manual work:

- dashboards
- formal workstation management screens
- scheduling UI
- machine/integration UI
- quality-control UI
- warehouse UI
- advanced analytics
