
-- connect as appuser

CREATE TABLE employees (
  id INT AUTO_INCREMENT PRIMARY KEY,
  first_name VARCHAR(25),
  last_name VARCHAR(25),
  email VARCHAR(75) NOT NULL,
  designation VARCHAR(50),
  date_of_joining DATE
) ;

-- sample data
INSERT INTO employees (first_name, last_name, email, designation, date_of_joining)
  VALUES
  ('David', 'Webb', 'david.webb@example.com', 'Principal', '2004-01-23'),
  ('Jason', 'Bourne', 'jason.bourne@example.com', 'Whoami', '2005-06-30'),
  ('John', 'Kane', 'john.kane@example.com', 'Staff', '2006-04-07'),
  ('Charles', 'Briggs', 'charles.briggs@example.com', 'Senior', '2007-11-05')
;

-- select * from employees;
