CREATE DATABASE bagipangan_user_test;
CREATE DATABASE bagipangan_request_test;
CREATE DATABASE bagipangan_food_test;
CREATE DATABASE bagipangan_claim_test;

\connect bagipangan_user_test
\i /migrations/user/001_init.sql
\connect bagipangan_request_test
\i /migrations/request/001_init.sql
\connect bagipangan_food_test
\i /migrations/food/001_init.sql
\connect bagipangan_claim_test
\i /migrations/claim/001_init.sql
