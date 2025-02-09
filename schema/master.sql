INSERT INTO user_cms (id, created_at, updated_at, deleted_at, username, password, flag, level, token) VALUES (1, null, 1719513238, null, 'admin', '21232f297a57a5a743894a0e4a801fc3', true, 'admin', '21232f297a57a5a743894a0e4a801fc3');
INSERT INTO user_cms (id, created_at, updated_at, deleted_at, username, password, flag, level, token) VALUES (2, null, 1718505303, null, 'report', 'e98d2f001da5678b39482efbdf5770dc', true, 'report', 'e98d2f001da5678b39482efbdf5770dc');

INSERT INTO programs (id, created_at, updated_at, deleted_at, retail, status, start_date, end_date) VALUES (1, '2024-06-01 22:17:26.955000 +00:00', '2024-06-15 15:08:35.760384 +00:00', null, 'am', true, 1718470800, 1722445199);

create type prize as enum ('gopay 25rb', 'abb', 'helm', 'jacket', 'gopay 50rb', 'gopay 100rb', 'handphone', 'smart tv', 'motor');
alter table public.prizes alter column prize type prize using prize::prize;

create type prize_type as enum ('reguler', 'zonk');
alter table public.prizes alter column prize_type type prize_type using prize_type::prize_type;

create type level_type as enum ('admin', 'report');
alter table public.user_cms alter column level type level_type using level::level_type;

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