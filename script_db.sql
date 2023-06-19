-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Hôte : 127.0.0.1:3306
-- Généré le : lun. 19 juin 2023 à 21:13
-- Version du serveur : 8.0.31
-- Version de PHP : 8.0.26

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Base de données : `forum`
--

-- --------------------------------------------------------

--
-- Structure de la table `categories`
--

DROP TABLE IF EXISTS `categories`;
CREATE TABLE IF NOT EXISTS `categories` (
  `id_category` int NOT NULL AUTO_INCREMENT,
  `category_title` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id_category`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb3;

--
-- Déchargement des données de la table `categories`
--

INSERT INTO `categories` (`id_category`, `category_title`) VALUES
(1, 'Voltaire'),
(2, 'Réseaux'),
(3, 'Challenge JS'),
(4, 'Forum'),
(5, 'Groupie Tracker'),
(6, 'Hangman'),
(7, 'Hangman Web'),
(8, 'Infra'),
(9, 'POO'),
(10, 'Linux'),
(11, 'Java'),
(12, 'Administration Poste Client'),
(13, 'Challenge 48H'),
(14, 'Ymmersion');

-- --------------------------------------------------------

--
-- Structure de la table `likers`
--

DROP TABLE IF EXISTS `likers`;
CREATE TABLE IF NOT EXISTS `likers` (
  `id_message` int NOT NULL,
  `id_user` int NOT NULL,
  PRIMARY KEY (`id_message`,`id_user`),
  KEY `id_user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

--
-- Déchargement des données de la table `likers`
--

INSERT INTO `likers` (`id_message`, `id_user`) VALUES
(2, 1),
(1, 2),
(1, 3),
(1, 4),
(2, 4),
(2, 5);

-- --------------------------------------------------------

--
-- Structure de la table `messages`
--

DROP TABLE IF EXISTS `messages`;
CREATE TABLE IF NOT EXISTS `messages` (
  `id_message` int NOT NULL AUTO_INCREMENT,
  `id_user` int DEFAULT NULL,
  `id_topic` int DEFAULT NULL,
  `content` text,
  `date_created` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_message`),
  KEY `id_user` (`id_user`),
  KEY `id_topic` (`id_topic`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb3;

--
-- Déchargement des données de la table `messages`
--

INSERT INTO `messages` (`id_message`, `id_user`, `id_topic`, `content`, `date_created`) VALUES
(1, 3, 1, 'Je sais pas, t\'es peut être un peu nul', '2023-06-10 15:11:05'),
(2, 1, 1, 'Je suis d\'accord avec toi !', '2023-06-10 15:15:00'),
(3, 2, 1, 'Oui, c\'est un vrai défi !', '2023-06-10 15:16:00');

-- --------------------------------------------------------

--
-- Structure de la table `reponses`
--

DROP TABLE IF EXISTS `reponses`;
CREATE TABLE IF NOT EXISTS `reponses` (
  `id_reponse` int NOT NULL AUTO_INCREMENT,
  `id_message` int DEFAULT NULL,
  `id_user` int DEFAULT NULL,
  `content` text,
  PRIMARY KEY (`id_reponse`),
  KEY `id_message` (`id_message`),
  KEY `id_user` (`id_user`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Déchargement des données de la table `reponses`
--

INSERT INTO `reponses` (`id_reponse`, `id_message`, `id_user`, `content`) VALUES
(2, 1, 2, 'Je suis d\'accord avec toi.'),
(3, 1, 3, 'C\'est pas très gentil'),
(4, 1, 1, 'Venant de lui ça m\'étonne pas...'),
(5, 2, 5, 'Je suis d\'accord avec celui qui est d\'accord avec toi');

-- --------------------------------------------------------

--
-- Structure de la table `roles`
--

DROP TABLE IF EXISTS `roles`;
CREATE TABLE IF NOT EXISTS `roles` (
  `id_role` int NOT NULL AUTO_INCREMENT,
  `name_role` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id_role`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb3;

--
-- Déchargement des données de la table `roles`
--

INSERT INTO `roles` (`id_role`, `name_role`) VALUES
(1, 'admin'),
(2, 'moderator'),
(3, 'user');

-- --------------------------------------------------------

--
-- Structure de la table `topics`
--

DROP TABLE IF EXISTS `topics`;
CREATE TABLE IF NOT EXISTS `topics` (
  `id_topic` int NOT NULL AUTO_INCREMENT,
  `id_category` int DEFAULT NULL,
  `id_user` int DEFAULT NULL,
  `topic_title` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id_topic`),
  KEY `id_category` (`id_category`),
  KEY `id_user` (`id_user`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb3;

--
-- Déchargement des données de la table `topics`
--

INSERT INTO `topics` (`id_topic`, `id_category`, `id_user`, `topic_title`) VALUES
(1, 1, 1, 'Voltaire c\'est trop dur'),
(2, 1, 2, 'Pourquoi Voltaire existe ???'),
(3, 5, 3, 'I need help !');

-- --------------------------------------------------------

--
-- Structure de la table `users`
--

DROP TABLE IF EXISTS `users`;
CREATE TABLE IF NOT EXISTS `users` (
  `id_user` int NOT NULL AUTO_INCREMENT,
  `id_role` int DEFAULT 3,
  `username` varchar(50) DEFAULT NULL,
  `email` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id_user`),
  KEY `id_role` (`id_role`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb3;

--
-- Déchargement des données de la table `users`
--

INSERT INTO `users` (`id_user`, `id_role`, `username`, `email`, `password`) VALUES
(1, 1, 'Charlemagne', 'test-valentin@gmail.com', 'azertyuiop'),
(2, 2, 'Plumecocq', 'test-augustin@gmail.com', 'azertyuiop'),
(3, 3, 'John', 'test-John@gmail.com', 'azertyuiop'),
(4, 3, 'Emma', 'emma@example.com', 'password123'),
(5, 3, 'Liam', 'liam@example.com', 'password456'),
(6, 1, 'admin', 'administrateur-forum@gmail.com', 'P4AC%85sxx'),
(7, NULL, 'fsqfsqfsfdq', 'fsqfsqf@', 'fqsfsqfsqfsq');

--
-- Contraintes pour les tables déchargées
--

--
-- Contraintes pour la table `likers`
--
ALTER TABLE `likers`
  ADD CONSTRAINT `likers_ibfk_1` FOREIGN KEY (`id_message`) REFERENCES `messages` (`id_message`),
  ADD CONSTRAINT `likers_ibfk_2` FOREIGN KEY (`id_user`) REFERENCES `users` (`id_user`);

--
-- Contraintes pour la table `messages`
--
ALTER TABLE `messages`
  ADD CONSTRAINT `messages_ibfk_1` FOREIGN KEY (`id_topic`) REFERENCES `topics` (`id_topic`),
  ADD CONSTRAINT `messages_ibfk_2` FOREIGN KEY (`id_user`) REFERENCES `users` (`id_user`);

--
-- Contraintes pour la table `topics`
--
ALTER TABLE `topics`
  ADD CONSTRAINT `fk_topics_categories` FOREIGN KEY (`id_category`) REFERENCES `categories` (`id_category`),
  ADD CONSTRAINT `fk_topics_users` FOREIGN KEY (`id_user`) REFERENCES `users` (`id_user`);

--
-- Contraintes pour la table `users`
--
ALTER TABLE `users`
  ADD CONSTRAINT `users_ibfk_1` FOREIGN KEY (`id_role`) REFERENCES `roles` (`id_role`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
