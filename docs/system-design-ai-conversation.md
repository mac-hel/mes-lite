## Question
Let's identify what systems are needed. I recon:

1. Backend API - registering production, managing orders, preparing reports etc.
2. Worker UI (should run everywhere: terminal/tablet/desktop, each employee or workstation has one) - this is worker's interface to insert production operations, number of manufactured products, encountered problems etc. Should be as simple, easy to use and ergonomic as possible.
3. UI for Manager/Leader/Admin (probably desktop) - provides management tools, insights, reports, admin panels etc.

I suppose 2. and 3. could be single application with panels/tools availability controlled by permissions.

## Answer
Yes — I’d treat this as **three logical systems, but probably only two deployable applications at first**:

1. **Backend / application server**
   Owns business logic and data. It handles production registration, orders, routing/operations, users and permissions, machine/workstation state, downtime/problems, quality data, reports, integrations, and audit history. For our current direction: **Go + PostgreSQL**, with HTTP endpoints and server-rendered views where appropriate.

2. **Unified MES web application**
   One browser application for both workers and management. Permissions and workstation context determine what is visible. This avoids maintaining two frontend codebases and makes cross-role workflows much simpler.

   The important distinction should be **different UX modes**, not different applications:

   * **Worker mode:** very large controls, minimal navigation, few decisions, fast repetitive input, touch-friendly, ideally usable with gloves/barcode scanner, almost no “administrative” clutter.
   * **Manager/Admin mode:** denser desktop UI with tables, filters, planning, reports, configuration, user management, corrections, audit views, etc.

3. **Integration / edge layer — optional initially, likely important later**
   I would explicitly reserve this as a separate subsystem even if version 1 barely uses it. Its job would be communication with machines, PLCs, barcode scanners, ERP, MQTT/OPC UA, printers, local devices, or other factory systems. Keeping this boundary separate prevents industrial connectivity concerns from contaminating the main application.

So architecturally I’d picture:

```text
Worker terminal ─┐
Tablet          ─┤
Desktop         ─┼──► MES Web Application
Manager PC      ─┤          │
Admin PC        ─┘          ▼
                         Go Backend
                             │
                        PostgreSQL
                             │
                 ┌───────────┴──────────┐
                 │                      │
            Reports/analytics      Integration/Edge
                                    ERP / PLC / machines
```

I agree strongly with combining **2 and 3**. Separate applications would only make sense if the worker interface later becomes something fundamentally different — for example an offline-first native kiosk app or a machine-specific HMI. Otherwise, a single application with **RBAC** (role-based access control) plus workstation-specific configuration is cleaner.

One nuance: permissions alone are not enough. I’d model both **user role** and **device/workstation context**. A production terminal may be locked to “Laser Cutting Station 3,” for example, so even a qualified worker sees only the actions relevant to that station. A manager logging into the same system from a desktop gets a completely different navigation structure.

That gives us a useful initial system breakdown: **Core MES backend, Unified Web UI, Database, and an optional Integration/Edge service boundary.** The next useful step would be to break those into concrete modules and decide which ones belong in the MVP.

## User remark
Initially **Unified MES web application** can be part of **Backend / application server**, just available on different URLs set which will render UI.
