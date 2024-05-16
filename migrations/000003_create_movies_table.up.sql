USE gosnip;

-- Add some dummy records (which we'll use in the next couple of chapters).
INSERT INTO snippets (title, content, created, expires) VALUES (
 'Dawn of Tech',
 'In the quiet hours of dawn, the world awakens to a symphony of code.\n Lines of logic intertwine, crafting a digital tapestry that powers the pulse of progress.',
 UTC_TIMESTAMP(),
 DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);
INSERT INTO snippets (title, content, created, expires) VALUES (
 'Whispers of Nature',
 'Beneath the emerald canopy, whispers of nature speak in rustling leaves.\n Each breath of wind carries stories from the heart of the earth,\n echoing the interconnected dance of life.',
 UTC_TIMESTAMP(),
 DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);
INSERT INTO snippets (title, content, created, expires) VALUES (
 'Enigmas of the Cosmos',
 'The cosmos is a vast canvas of enigmas, where stars are born from cosmic dust\n and black holes weave the fabric of spacetime into a dance of gravity.\n Each galaxy holds a billion mysteries, whispering the secrets of\n the universe in the language of light.',
 UTC_TIMESTAMP(),
 DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);
