INSERT INTO user_cms (id, created_at, updated_at, deleted_at, username, password, flag, level, token) VALUES (1, null, 1719513238, null, 'admin', '3515695a8e1e1e97a8af6631e6de23ae', true, 'admin', 'e6342ea745515132aabb697bb218cd7a');
INSERT INTO user_cms (id, created_at, updated_at, deleted_at, username, password, flag, level, token) VALUES (2, null, 1718505303, null, 'report', '7d5ba8421fce3535a1ff747faf84a4e4', true, 'report', '9f7757618fb5e25847af40d377d2a68d');

INSERT INTO programs (id, created_at, updated_at, deleted_at, retail, status, start_date, end_date) VALUES (1, '2024-06-01 22:17:26.955000 +00:00', '2024-06-15 15:08:35.760384 +00:00', null, 'am', true, 1718470800, 1722445199);

create type prize as enum ('gopay 25rb', 'abb', 'helm', 'jacket');
alter table public.prizes alter column prize type prize using prize::prize;

create type prize_type as enum ('reguler', 'zonk');
alter table public.prizes alter column prize_type type prize_type using prize_type::prize_type;

create view detailed_prize_redemptions
            (prize, name, nik, msisdn, date_redeem,is_zonk, redeem_id, county, profession, date_validation, approved, lottery_number, amount)
as
SELECT p.prize,
       u.name,
       u.nik,
       u.wa_id                  AS msisdn,
       redeem_prizes.created_at AS date_redeem,
       u.is_zonk,
       redeem_prizes.id         AS redeem_id,
       u.county,
       u.profession,
       redeem_prizes.date_validation,
       redeem_prizes.approved,
       redeem_prizes.lottery_number,
       redeem_prizes.amount,
       redeem_prizes.notes
FROM redeem_prizes
         JOIN prizes p ON redeem_prizes.prize_id = p.id
         JOIN users_unique_codes u ON redeem_prizes.users_unique_code_id = u.id
ORDER BY redeem_prizes.created_at DESC;