-- 1. Добавляем временный столбец для одной строки
ALTER TABLE news ADD COLUMN image_url_tmp TEXT;

-- 2. Копируем из массива первый элемент обратно в строку
UPDATE news SET image_url_tmp = image_urls[1];

-- 3. Удаляем столбец с массивом
ALTER TABLE news DROP COLUMN image_urls;

-- 4. Переименовываем временный столбец обратно
ALTER TABLE news RENAME COLUMN image_url_tmp TO image_url;
