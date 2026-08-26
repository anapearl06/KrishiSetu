# 🌾 KrishiSetu

> **Connecting Farmers Directly to Market Demand**

[![SIH 2026](https://img.shields.io/badge/SIH-26033-orange.svg)](https://sih.gov.in)
[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)
[![FastAPI](https://img.shields.io/badge/FastAPI-0.100+-green.svg)](https://fastapi.tiangolo.com)
[![Flutter](https://img.shields.io/badge/Flutter-3.0+-02569B.svg)](https://flutter.dev)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-4169E1.svg)](https://www.postgresql.org)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**KrishiSetu** is an AI-powered agricultural marketplace developed for **SIH 26033** (Ministry of Consumer Affairs, Food & Public Distribution, Department of Consumer Affairs). It connects farmers and FPOs (Farmer Producer Organizations) directly with consumers and bulk buyers, leveraging agricultural data, demand forecasting, price intelligence, and intelligent matching to reduce supply-demand inefficiencies and eliminate middlemen exploitation.

---

## 🎯 Problem Statement

The traditional agricultural supply chain suffers from multiple structural inefficiencies:
* **Intermediary Margins:** Multiple layers of middlemen significantly reduce farmer earnings while inflating consumer prices.
* **Information Asymmetry:** Farmers have limited visibility into real-time market demand and pricing trends across different mandis.
* **Buyer Discovery Friction:** Bulk buyers and food processors struggle to reliably source verified produce matching specific quantity and quality criteria.
* **Price Volatility:** Severe price fluctuations across regional markets lead to distress sales and food waste.
* **Supply-Demand Mismatches:** Lack of predictive analytics results in local oversupply and shortages.

---

## 💡 Solution

KrishiSetu acts as a comprehensive digital bridge between **farm supply and market demand**, driven by data intelligence.

```text
Farmer Supply ───────┐
                     ▼
              ┌─────────────┐
Buyer Demand ─► KrishiSetu  │
              │ Intelligence│
              └──────┬──────┘
                     │
          ┌──────────┼──────────┐
          ▼          ▼          ▼
       Demand      Price      Smart
      Forecasting Intelligence Matching
          └──────────┼──────────┘
                     ▼
                  Order
```

---

## 🚀 Core Features

### 👨‍🌾 Farmer / FPO Module
* **Registration & Profile Verification:** Secure onboarding with land/FPO verification status.
* **Produce Listings:** Manage crop variety, quantity, quality grading, expected price, and location.
* **Availability Scheduling:** Set harvest dates and active availability windows.
* **Market Insights:** Real-time access to regional mandi prices and trends.
* **Buyer Recommendations:** AI-driven buyer matching for listed produce.
* **Offer & Order Management:** Receive, negotiate, accept, or reject purchase offers.
* **Payment Status Tracking:** Transparent tracking of transaction and payment milestones.

### 🏪 Buyer Module
* **Registration & Profile Management:** Tailored onboarding for retailers, wholesalers, and food processors.
* **Demand Creation:** Post detailed crop requirements including quantity, quality specs, budget, and required dates.
* **Smart Recommendations:** Algorithmic ranking of suitable farmers based on proximity, quality, and price.
* **Farmer Comparison:** Side-by-side comparison of candidate farmers.
* **Order Tracking:** End-to-end visibility from offer confirmation to fulfillment.

### 🧠 Intelligence Module
* **Demand Forecasting:** Predict future crop demand using historical consumption, seasonality, and market arrivals.
* **Price Intelligence:** Historical price trends, regional comparisons, and reference pricing.
* **Explainable Matching:** Transparent scoring system detailing *why* a particular farmer matches a buyer's demand.
* **Regional Market Analysis:** Aggregated state/district-level agricultural market insights.
* **Distance & Transportation Estimation:** Calculate travel distance, transit time, and estimated logistics costs using OSRM.

---

## 🔄 Complete Flow

```text
                    FARMER
                      │
                Lists Produce
                      │
                      ▼
                ┌──────────┐
                │          │
                │KrishiSetu│◄──── Buyer creates Demand
                │          │
                └────┬─────┘
                     │
          ┌──────────┼──────────┐
          ▼          ▼          ▼
       Market      Demand     Matching
        Data      Forecast    Engine
          │          │          │
          └──────────┼──────────┘
                     ▼
              Ranked Farmers
                     │
                     ▼
                  Offer
                     │
                     ▼
                  Order
                     │
                     ▼
              Payment Status
```

### Farmer Workflow
1. Register & Create Profile
2. List Produce (Crop, Variety, Quantity, Quality)
3. Set Availability & Expected Price
4. View Market Insights & Buyer Opportunities
5. Receive & Review Buyer Offers (Accept / Reject)
6. Fulfill Order & Track Payment Status

### Buyer Workflow
1. Register & Create Profile
2. Create Demand (Crop, Quantity, Quality, Budget, Location)
3. Submit Demand to Matching Engine
4. Review Ranked Recommendations with Explainability Scores
5. Select Farmer & Send Offer / Order
6. Track Order Fulfillment & Payment Status

---

## 🎯 Intelligent Matching Algorithm

The core matching engine scores and ranks farmers against buyer demand using weighted multi-criteria decision analysis:

$$	ext{Match Score} = w_1(	ext{Crop}) + w_2(	ext{Quantity}) + w_3(	ext{Price}) + w_4(	ext{Quality}) + w_5(	ext{Availability}) + w_6(	ext{Distance}) + w_7(	ext{Market Signals})$$

### Example Breakdown: **Farmer A — 94% Match**
* ✓ Required crop and variety available
* ✓ Quantity requirement fully satisfied
* ✓ Expected price within buyer budget range
* ✓ Required quality grade met
* ✓ Suitable geographic distance (optimized transport cost)
* ✓ Available before required fulfillment date

> **Transparent & Explainable:** The system prioritizes explainable recommendations with clear criteria breakdowns instead of opaque black-box scores.

---

## 🧠 AI / ML Architecture

### Demand Forecasting Pipeline
```text
Agricultural Data ──> Data Cleaning ──> Feature Engineering ──> Baseline Models ──> ML Models (XGBoost / Random Forest) ──> Evaluation (MAE, RMSE, MAPE) ──> Demand Prediction
```

* **Potential Models:** Linear Regression, Random Forest, XGBoost, and Time-series models.
* **Features Used:** Historical demand, historical prices, market arrivals, crop production statistics, region, month, season, lag features, and rolling averages.
* **Validation:** Temporal validation and rigorous data leakage checks.

### Price Intelligence
Leverages historical agricultural market data to provide:
* Price trends & movement analysis
* Regional and mandi-wise price comparisons
* Market reference prices & volatility metrics

---

## 🏗️ Architecture & System Design

```text
                       ┌───────────────┐
                       │    Flutter    │
                       │   Frontend    │
                       └───────┬───────┘
                               │
                               ▼
                       ┌───────────────┐
                       │   Go + Gin    │
                       │   REST API    │
                       └───────┬───────┘
                               │
             ┌─────────────────┼─────────────────┐
             │                 │                 │
             ▼                 ▼                 ▼
      ┌─────────────┐   ┌─────────────┐   ┌──────────────┐
      │ PostgreSQL  │   │    Redis    │   │ Python API   │
      │             │   │             │   │   FastAPI    │
      └─────────────┘   └─────────────┘   └──────┬───────┘
                                                  │
                                                  ▼
                                        ┌──────────────────┐
                                        │ ML / Intelligence│
                                        │                  │
                                        │ Demand Forecast  │
                                        │ Price Analysis   │
                                        │ Matching Signals │
                                        └──────────────────┘
```

### Backend Request Lifecycle
```text
HTTP Request ──> Middleware ──> Handler ──> Service ──> Repository ──> GORM ──> PostgreSQL
```
* **Go + Gin:** Handles primary application services, RESTful routing, authentication, and core business logic.
* **Python + FastAPI:** Isolated data-science and machine learning inference layer.

---

## ⚙️ Tech Stack

| Layer | Technology |
| :--- | :--- |
| **Frontend** | Flutter (Cross-platform Mobile & Web) |
| **Backend API** | Go, Gin Framework |
| **ORM & Database** | GORM, PostgreSQL 15+ |
| **Caching & Session** | Redis |
| **ML & Analytics** | Python, Pandas, NumPy, Scikit-learn, XGBoost |
| **ML API** | FastAPI |
| **Maps & Routing** | OpenStreetMap, OSRM |
| **Security** | JWT, bcrypt, Role-Based Access Control (RBAC) |
| **Infrastructure** | Docker, Docker Compose, Nginx |
| **Documentation & Testing** | Swagger / OpenAPI, Go Testing & Testify |

---

## 🔌 Request Flows

### 1. User Login
```text
Flutter ──> POST /api/v1/auth/login ──> Gin Router ──> Auth Handler ──> Auth Service ──> User Repository ──> PostgreSQL ──> Password Verification ──> JWT Generation ──> Response
```

### 2. Create Listing
```text
Flutter ──> POST /api/v1/listings ──> JWT Middleware ──> Listing Handler ──> Listing Service ──> Validation ──> Repository ──> PostgreSQL ──> Response
```

### 3. Intelligent Matching & ML Prediction
```text
Buyer Demand ──> Go API ──> Matching Service ──> Candidate Farmers (PostgreSQL) ──> FastAPI (ML Signals) ──> Scoring & Ranking ──> Explainable Recommendations ──> Flutter
```

---

## 🗄️ Core Entities & Data Model

### Key Entities
`User`, `Farmer`, `FPO`, `Buyer`, `Produce`, `Listing`, `Demand`, `Match`, `Offer`, `Order`, `Inventory`, `Payment`, `Notification`, `Analytics`.

### Schema Highlights

* **Farmer / FPO:** `id`, `user_id`, `name`, `location`, `verification_status`
* **Listing:** `id`, `farmer_id`, `crop`, `variety`, `quantity`, `quality`, `expected_price`, `location`, `harvest_date`, `available_from`, `status`
* **Buyer Demand:** `id`, `buyer_id`, `crop`, `quantity`, `quality_requirement`, `max_price`, `location`, `required_date`, `status`
* **Match:** `id`, `demand_id`, `farmer_id`, `score`, `distance`, `price_score`, `quantity_score`, `quality_score`, `availability_score`, `explanation`
* **Order:** `id`, `buyer_id`, `farmer_id`, `listing_id`, `quantity`, `price`, `total_amount`, `status`, `payment_status`, `created_at`, `updated_at`

---

## 📊 Data Sources

KrishiSetu utilizes publicly available agricultural data from official Indian repositories:
* **Government of India Open Data Platform:** [data.gov.in](https://data.gov.in)
* **Agmarknet:** Mandi arrivals and pricing datasets.
* **Crop Production Datasets:** State and district-level yield, area, and production metrics.
* **Geographic Reference Datasets:** State/district boundary and routing data.

> **Mandi Data Fields:** `Date`, `State`, `District`, `Market`, `Commodity`, `Variety`, `Arrival Quantity`, `Minimum Price`, `Maximum Price`, `Modal Price`.  
> **Crop Production Data Fields:** `Year`, `State`, `District`, `Crop`, `Area`, `Production`, `Yield`, `Season`.  
> *Note: Where real-time transaction data is unavailable during development, clearly labelled synthetic data is used for demonstration.*

---

## 🔐 Security, Reliability & Inventory Safety

* **Authentication & Authorization:** Secure JWT tokens, bcrypt password hashing, and granular Role-Based Access Control (RBAC).
* **API Protection:** Request validation middleware, CORS, and robust error handling.

### 🛡️ Inventory Safety & Concurrency Control
To prevent race conditions and overselling:
```text
Available Inventory = 10 tonnes
├── Buyer A requests 7 tonnes (Approved)
└── Buyer B requests 5 tonnes (Blocked / Partial Allocation)
```
* **PostgreSQL Transactions:** Database-level locking and ACID transactions ensure absolute inventory consistency and prevent overselling during high-concurrency demand spikes.

---

## 🚚 Logistics & Payment Scope

### Logistics Scope (Intelligence Only)
* **Included:** Farmer-buyer distance calculation, estimated travel distance & time, approximate transportation cost estimation, and basic route guidance via OSRM.
* **Excluded:** Driver management, fleet management, vehicle assignment, live vehicle tracking, warehouse management, and physical transport operations.

### Payment Scope (Status Tracking)
* **Included:** Transaction milestone tracking (`PENDING` ➔ `INITIATED` ➔ `COMPLETED`).
* **Excluded:** Full payment gateway infrastructure and final settlement escrow (designed to integrate with external payment gateways like Razorpay/UPI in production).

---

## 🔒 MVP Scope

### ✅ Included in MVP
* Farmer/FPO & Buyer marketplaces
* Produce listings and buyer demand posting
* Intelligent multi-criteria matching engine
* Demand forecasting and price intelligence dashboards
* Regional market insights and distance/cost estimation
* Offer negotiation, order management, and inventory safety
* Payment status tracking and notifications

### ❌ Excluded from MVP (Post-MVP Roadmap)
* Driver and fleet management
* Physical logistics and warehousing
* Live vehicle tracking
* Full payment gateway settlement engines
* IoT sensor integration and Blockchain tracking
* Unnecessary microservices overhead

---

## 🌍 Impact

KrishiSetu delivers measurable socio-economic benefits:
* **Enhanced Farmer Earnings:** Direct market access cuts out intermediaries, securing better margins for farmers.
* **Price Transparency:** Real-time mandi price intelligence prevents distress selling.
* **Demand Visibility:** Predictive analytics guide farmers on what crops to grow and when to harvest.
* **Efficient Matching:** Drastically reduces discovery friction between bulk buyers and agrarian producers.
* **Market Efficiency:** Stabilizes supply chains and mitigates local price shocks.

---

## 🏆 Smart India Hackathon 2026 (SIH)

* **Problem Statement ID:** `26033`
* **Organization:** Ministry of Consumer Affairs, Food & Public Distribution
* **Department:** Department of Consumer Affairs
* **Category:** Software

---

<p align="center">
  <b>🌾 KrishiSetu</b><br>
  <i>From Farm Supply to Market Demand.</i>
</p>
