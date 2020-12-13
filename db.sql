USE servente;

CREATE TABLE `services` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL,
  `content` text NOT NULL,
  `created` datetime NOT NULL,
  `expires` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_services_created` (`created`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=latin1;

CREATE TABLE `users2` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `hashed_password` char(60) NOT NULL,
  `created` datetime NOT NULL,
  `active` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_uc_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
