# Library music
# Пример использования API

## API Endpoint: Create
Endpoint для добавления новой песни.

### Request
- Method: `POST`
- URL: `http://localhost:8080/api/v1/song`
- Headers:
  - `Content-Type: application/json`
- Body:
  ```json
  {
      "group": "Lady Gaga",
      "name": "Poker Face",
      "link": "https://genius.com/Lady-gaga-poker-face-lyrics",
      "text": [                  
          {
              "type": "verse",
              "text": "I wanna hold 'em like they do in Texas, please (Woo)"
          },
          {
              "type": "verse",
              "text": "Fold 'em, let 'em hit me, raise it, baby, stay with me (I love it)"
          },
          {
              "type": "verse",
              "text": "LoveGame intuition, play the cards with spades to start"
          },
          {
              "type": "verse",
              "text": "And after he's been hooked, I'll play the one that's on his heart"
          },
          {
              "type": "chorus",
              "text": "Can't read my, can't read my"
          },     
          {
              "type": "chorus",
              "text": "No, he can't read my poker face"
          }
      ],
      "release_date": "2008-09-23"
  }
  ```

### Response
- **Success Response:**
  - Code: `200`
  - Body:
    ```json
    {
        "response": "id"
    }
    ```
- **Incorrect data:**
  - Code: `400`
  - Body:
    ```json
    {
        "error": "string"
    }
    ```
- **InternalServerError:**
  - Code: `500`
  - Body:
    ```json
    {
        "error": "string"
    }
    ```
    
## API Endpoint: Update
Endpoint для обновления песни.

### Request
- Method: `PATCH`
- URL: `http://localhost:8080/api/v1/song/{id}`
- Headers:
  - `Content-Type: application/json`
- Body:
  ```json
  {
      "group": "Lady Gaga",
      "name": "Poker Face",
      "link": "https://lyrsense.com/lady_gaga/poker_face",
      "release_date": "2008-09-23",
      "text": [{
          "type": "verse",
          "text": "I wanna hold 'em like they do in Texas, please (Woo)",
      }]
  }
  ```

### Response
- **Success Response:**
  - Code: `200`  
  - Body:
    ```json
    {
        "response": "ok"
    }
    ```
- **Incorrect data:**
  - Code: `400`
  - Body:
    ```json
    {
        "error": "string"
    }
    ```
- **Not Found:**
  - Code: `404`  
  - Body:
    ```json
    {
        "error": "song not found"
    }
    ```
- **InternalServerError:**
  - Code: `500`
  - Body:
    ```json
    {
        "error": "string"
    }
    ```

## API Endpoint: Delete
Endpoint для удаления песни.

### Request
- Method: `Delete`
- URL: `http://localhost:8080/api/v1/song/{id}`
- Headers:
  - `Content-Type: application/json`

### Response
- **Success Response:**
  - Code: `200`  
  - Body:
    ```json
    {
        "response": "ok"
    }
    ```
- **Incorrect data:**
  - Code: `400`
  - Body:
    ```json
    {
        "error": "string"
    }
    ```
- **Not Found:**
  - Code: `404`  
  - Body:
    ```json
    {
        "error": "song not found"
    }
    ```
- **InternalServerError:**
  - Code: `500`
  - Body:
    ```json
    {
        "error": "string"
    }
    ```

## API Endpoint: GetSong
Endpoint для получения текста песни с пагинацией по куплетам.

### Request
- Method: `Get`
- URL: `http://localhost:8080/api/v1/song`
- Headers:
  - `Content-Type: application/json`
- Params:
  - `group: "Lady Gaga"`
  - `name: "Poker Face"`
  - `offset: 1`

### Response
- **Success Response:**
  - Code: `200`
  - Body:
    ```json
      {
          "response": "I wanna hold 'em like they do in Texas, please (Woo)"
      }
    ```
- **Incorrect data:**
  - Code: `400`
  - Body:
    ```json
    {
        "error": "string"
    }
    ```
- **Not Found:**
  - Code: `404`  
  - Body:
    ```json
    {
        "error": "song not found"
    }
    ```
- **InternalServerError:**
  - Code: `500`
  - Body:
    ```json
    {
        "error": "string"
    }
    ```

## API Endpoint: GetSongs
Endpoint для получения данных библиотеки с фильтрацией по всем полям и пагинацией

### Request
- Method: `Get`
- URL: `http://localhost:8080/api/v1/songs`
- Headers:
  - `Content-Type: application/json`
- Params:
  - `group: "L"`
  - `name: "P"`
  - `link: "https://lyrsense.com"`
  - `release_date: 2006-01-01`
  - `offset: 0`
  - `limit: 2`

### Response
- **Success Response:**
  - Code: `200`
  - Body:
    ```json
      {
          "response": [{
              "group": "Lady Gaga",
              "name": "Poker Face",
              "link": "https://lyrsense.com/lady_gaga/poker_face",
              "release_date": "2008-09-23",
          }]      
      }
    ```
- **Incorrect data:**
  - Code: `400`
  - Body:
    ```json
    {
        "error": "string"
    }
    ```
- **Not Found:**
  - Code: `404`  
  - Body:
    ```json
    {
        "error": "song not found"
    }
    ```
- **InternalServerError:**
  - Code: `500`
  - Body:
    ```json
    {
        "error": "string"
    }
    ```