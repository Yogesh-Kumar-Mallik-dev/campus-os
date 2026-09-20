# ADR-0004: Academic Course-Semester Structure & Department Hosting Model

* **Status:** Accepted / Locked
* **Date:** 2026-09-20
* **Deciders:** Developer, Yogesh

---

## Context

The relationship between Programs, Courses, Departments, Semesters, and operational entities (Sections, Timetables, Teachers) required unambiguous institutional modeling.

## Decision

1. **Curriculum Hierarchy:**
   * A **Course** (e.g. B.Tech Computer Science) represents the degree curriculum.
   * A Course contains sequentially structured **Semesters** (e.g. 1st through 8th Semester).
   * A Semester contains defined **Subjects** and **Labs**.
2. **Department Hosting Model:**
   * An academic **Department** possesses the institutional authority to host specific semesters for one or more courses.
   * Within a hosted semester instance, the Department has full operational power to:
     * Allocate **Sections** (e.g., Section A, Section B).
     * Configure **Timetables** and lecture/lab slots.
     * Assign **Teachers** and **Lab Technicians** to subjects and sections.

## Consequences

* **Positive:** Accurately mirrors collegiate reality where departments host multiple concurrent semesters across undergraduate and postgraduate programs while retaining operational autonomy over scheduling and staff allocation.
* **Negative:** Requires queries that correlate curriculum-level definitions (`academic_schema.courses`) with operational hosted instances (`academic_schema.hosted_semesters`).
