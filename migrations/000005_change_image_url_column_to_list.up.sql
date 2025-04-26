-- 1. Добавляем временный столбец
ALTER TABLE news ADD COLUMN image_urls_tmp TEXT[];

-- 2. Копируем данные, оборачивая старую строку в массив
UPDATE news SET image_urls_tmp = ARRAY[image_url];

-- 3. Удаляем старый столбец
ALTER TABLE news DROP COLUMN image_url;

-- 4. Переименовываем новый столбец
ALTER TABLE news RENAME COLUMN image_urls_tmp TO image_urls;
