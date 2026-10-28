# learNuggget

Eine einfache Flashcard-Web-App, mit der du eigene Kollektionen und Karten erstellen und mit verschiedenen Lernmodi üben kannst.

## Aktueller Funktionsumfang

- Kollektionen erstellen und löschen
- Einzelne Karten hinzufügen oder mehrere Karten auf einmal importieren
- Anzahl der Karten pro Kollektion anzeigen
- Kollektionen nach Neueste, Älteste, A–Z oder Z–A sortieren
- Karten während des Lernens löschen
- Mit drei Lernmodi üben:
    - **Infinity:** Karten unbegrenzt durchblättern, ohne Bewertung oder Ende.
    - **Practice:** Karten als richtig oder falsch markieren. Richtig beantwortete Karten werden aus der aktuellen Runde entfernt; falsch beantwortete Karten bleiben zum Weiterüben verfügbar. Die Runde endet, wenn alle Karten als richtig markiert wurden. Danach wird eine Auswertung angezeigt.
    - **One-try:** Jede Karte kann einmal beantwortet werden. Die Runde endet, wenn alle Karten bearbeitet wurden. Danach wird eine Auswertung angezeigt.
- In Infinity mit Tastatur und Wischgesten navigieren

> Scores und Fortschritt werden derzeit nur während der aktuellen Runde erfasst. Die App speichert noch keine Scores, Benutzerkonten oder Lernhistorie.

## Verwendete Technologien

- **Frontend:** Svelte 5 und Vite
- **Backend:** Go mit Gin
- **Datenbank:** SQLite

## Installation und Start

### Voraussetzungen

- Node.js und npm
- Go

### Backend starten

```bash
cd backend
go run .
```

Die API startet unter `http://localhost:3000`. Beim Start werden die benötigten SQLite-Tabellen erstellt. Die Datenbankdatei wird unter `backend/learNuggget.db` angelegt.

### Frontend starten

Öffne ein zweites Terminal im Projektverzeichnis und führe Folgendes aus:

```bash
npm install
npm run dev
```

Öffne anschließend die lokale Adresse, die Vite im Terminal anzeigt. Während der Entwicklung leitet Vite Anfragen unter `/api` an das Backend auf Port `3000` weiter.

Für einen Produktions-Build des Frontends:

```bash
npm run build
```

## Projektstruktur

```text
src/
  lib/          Svelte-Komponenten für Karten, Kollektionen und Auswertungen
  pages/        Startseite und Lernseite
  style/        Styles der Anwendung
backend/
  handlers/     HTTP-Handler
  routes/       Registrierung der API-Routen
  functions/    Datenbankzugriff und Anwendungslogik
```

## API-Übersicht

| Methode | Endpunkt | Beschreibung |
| --- | --- | --- |
| `GET` | `/collections?orderState=0` | Kollektionen abrufen; Sortierung: `0` neueste, `1` älteste, `2` A–Z, `3` Z–A |
| `POST` | `/collections/new?name=<name>` | Kollektion erstellen |
| `DELETE` | `/collections/delete?name=<name>` | Kollektion und ihre Karten löschen |
| `GET` | `/collection/len?id=<id>` | Anzahl der Karten einer Kollektion abrufen |
| `GET` | `/cards?collectionID=<id>` | Karten einer Kollektion abrufen |
| `POST` | `/addSingleCard` | Einzelne Karte hinzufügen |
| `POST` | `/addBulk?id=<id>` | Mehrere Karten auf einmal hinzufügen |
| `DELETE` | `/card/delete?id=<id>` | Karte löschen |

## Ideen für die Zukunft

Die folgenden Funktionen sind geplant, aber noch nicht implementiert.

### Benutzerkonten und soziale Funktionen

- Benutzerkonten
- Ein Profil-Gesamtscore, berechnet aus dem Durchschnitt aller Spielmodi und Kollektionen
- Den jeweils letzten Score pro Spielmodus und Kollektion speichern
- Freunde hinzufügen und Scores miteinander vergleichen

### Karten markieren

- Einzelne Karten in einer Kollektion markieren
- Nur markierte Karten üben
- Anzahl der markierten Karten anzeigen

### Spielmodi und Bewertung

- **Infinity:** Karten unbegrenzt durchlaufen, ohne Richtig-/Falsch-Bewertung oder Ende.
- **Practice:** Unbegrenzt viele Versuche ermöglichen. Die Runde endet, wenn alle Karten als richtig markiert wurden. Geplanter Score:
  `100 ÷ (1 + (falscheKlicks × 2) ÷ anzahlKarten)`
- **One-try:** Einen Versuch pro Karte erlauben, ohne Rotation oder Wiederholung. Geplanter Score:
  `(falscheAntworten ÷ anzahlKarten) × 100`

### Geteilte Kollektionen und Ranglisten

- Globale Kollektionen, die vom Root-Benutzer erstellt werden und von allen Benutzern gesehen und bearbeitet werden können
- Highscore und Benutzername für jede globale Kollektion anzeigen
- Globale Top 3 der Profile mit dem höchsten Profilscore

### Organisation und Darstellung

- Dark Mode
- Kollektionen in Gruppen organisieren
- Level-System

## Lizenz

Siehe [LICENSE](LICENSE).