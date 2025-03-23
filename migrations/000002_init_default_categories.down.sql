DELETE FROM categories 
WHERE name IN ('Еда', 'Транспорт', 'Коммуналка', 'Здоровье', 'Одежда', 'Развлечения', 'Спорт', 'Подарки', 'Прочее')
  AND is_default = true;
